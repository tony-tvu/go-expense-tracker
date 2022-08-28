package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/plaidapi"
	"github.com/tony-tvu/goexpense/user"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	App *app.App
}

var env string
var port string
var mongoURI string
var dbName string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}
	env = os.Getenv("ENV")
	port = os.Getenv("PORT")
	mongoURI = os.Getenv("MONGODB_URI")
	dbName = os.Getenv("DB_NAME")
	if util.ContainsEmpty(env, port, mongoURI, dbName) {
		log.Fatal("env variables are missing")
	}
}

func (s *Server) Initialize() {

	// Init handlers
	authHandler := &auth.AuthHandler{App: s.App}
	userHandler := &user.UserHandler{App: s.App}

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	// TODO: add logging middlware
	
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
		authRequired.GET("/create_link_token", plaidapi.CreateLinkToken)
		authRequired.POST("/set_access_token", plaidapi.SetAccessToken)

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
	if env == "development" || env == "production" {
		mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
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

		s.App.Users = mongoclient.Database(dbName).Collection("users")
		s.App.Sessions = mongoclient.Database(dbName).Collection("sessions")
	}

	h := cors.New(cors.Options{
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"Content-Type", "Plaid-Public-Token"},
		AllowCredentials: true,
	}).Handler(s.App.Router)

	srv := &http.Server{
		Handler:      h,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Listening on port %s", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
