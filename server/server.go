package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/plaid/plaid-go/plaid"
	"github.com/rs/cors"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/plaidapi"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	App *app.App
}

func (s *Server) Initialize() {
	if string(s.App.EncryptionKey) == "" {
		log.Fatal("fatal: missing ENCRYPTION_KEY")
	}
	if string(s.App.JwtKey) == "" {
		log.Fatal("fatal: missing JWT_KEY")
	}
	if s.App.MongoURI == "" {
		log.Fatal("fatal: missing MONGODB_URI")
	}
	if s.App.Env == "" {
		s.App.Env = "development"
	}
	if s.App.Port == "" {
		s.App.Port = "80"
	}

	// Plaid client
	if s.App.PlaidClientID != "" ||
		s.App.PlaidSecret != "" ||
		s.App.PlaidEnv != "" ||
		s.App.PlaidProducts != "" ||
		s.App.PlaidCountryCodes != "" {
		plaidCfg := plaid.NewConfiguration()
		plaidCfg.AddDefaultHeader("PLAID-CLIENT-ID", s.App.PlaidClientID)
		plaidCfg.AddDefaultHeader("PLAID-SECRET", s.App.PlaidSecret)
		plaidCfg.UseEnvironment(app.PlaidEnvs[s.App.PlaidEnv])
		plaidClient := plaid.NewAPIClient(plaidCfg)
		s.App.PlaidClient = plaidClient
	} else {
		log.Println("plaid configs are missing - service will not work")
	}

	// Init handlers
	authHandler := &auth.AuthHandler{App: s.App}
	plaidHandler := &plaidapi.PlaidHandler{App: s.App}
	userHandler := &user.UserHandler{App: s.App}

	if s.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.ForwardedByClientIP = true

	// apply global middleware
	router.Use(middleware.RateLimit())
	router.Use(middleware.NoCache)

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Ok"})
	})
	router.POST("/api/login", authHandler.Login, middleware.LoginRateLimit())

	authRequired := router.Group("/api", middleware.AuthRequired(s.App))
	{
		authRequired.POST("/logout", authHandler.Logout)
		authRequired.GET("/user_info", userHandler.GetInfo)
		authRequired.GET("/create_link_token", plaidHandler.CreateLinkToken)
		authRequired.POST("/set_access_token", plaidHandler.SetAccessToken)

		adminRequired := authRequired.Group("/", middleware.AdminRequired(s.App))
		{
			adminRequired.POST("/invite", userHandler.Invite)
			adminRequired.GET("/sessions", userHandler.GetSessions)
		}
	}

	// serve frontend
	router.Use(static.Serve("/", static.LocalFile("./web/build", true)))
	// prevent returning 404 when reloading page on frontend route
	router.NoRoute(func(ctx *gin.Context) {
		ctx.File("./web/build")
	})

	s.App.Router = router
}

func (s *Server) Run(ctx context.Context) {
	if s.App.Env != "test" {
		mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(s.App.MongoURI))
		if err != nil {
			log.Fatal(err)
		}
		err = mongoclient.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := mongoclient.Disconnect(ctx); err != nil {
				log.Println("mongo has been disconnected: ", err)
			}
		}()
		if s.App.DbName == "" {
			s.App.DbName = "goexpense_local"
		}
		s.App.Users = mongoclient.Database(s.App.DbName).Collection("users")
		s.App.Sessions = mongoclient.Database(s.App.DbName).Collection("sessions")
	}

	h := cors.New(cors.Options{
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"Content-Type", "Plaid-Public-Token"},
		AllowCredentials: true,
	}).Handler(s.App.Router)

	srv := &http.Server{
		Handler:      h,
		Addr:         fmt.Sprintf(":%s", s.App.Port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Listening on port %s", s.App.Port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
