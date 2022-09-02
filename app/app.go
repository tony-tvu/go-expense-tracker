package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/plaidapi"
	"github.com/tony-tvu/goexpense/user"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Db     *database.Db
	Router *gin.Engine
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

func (a *App) Initialize(ctx context.Context) {
	a.Db = &database.Db{}
	u := &user.UserHandler{Db: a.Db}

	if env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.ForwardedByClientIP = true
	if env == "development" {
		allowCrossOrigin(router)
	}

	// global middleware
	router.Use(middleware.RateLimit())
	router.Use(middleware.Logger(env))

	apiGroup := router.Group("/api", middleware.NoCache)
	{
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Ok"})
		})
		apiGroup.POST("/logout", u.Logout)
		apiGroup.POST("/login", middleware.LoginRateLimit(), u.Login)

		authGroup := apiGroup.Group("/", middleware.AuthRequired(a.Db))
		{
			authGroup.GET("/ping", u.Ping)
			authGroup.GET("/user_info", u.GetUserInfo)
			authGroup.GET("/create_link_token", plaidapi.CreateLinkToken)
			authGroup.POST("/set_access_token", plaidapi.SetAccessToken)

			adminGroup := authGroup.Group("/", middleware.AdminRequired(a.Db))
			{
				adminGroup.POST("/invite", u.InviteUser)
				adminGroup.GET("/sessions", u.GetSessions)
			}
		}

	}

	router.Use(middleware.FrontendCache, static.Serve("/", static.LocalFile("./web/build", true)))
	router.NoRoute(middleware.FrontendCache, func(ctx *gin.Context) {
		ctx.File("./web/build")
	})

	a.Router = router
}

func (a *App) Run(ctx context.Context) {
	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()
	a.Db.Users = mongoclient.Database(dbName).Collection("users")
	a.Db.Sessions = mongoclient.Database(dbName).Collection("sessions")
	createInitialAdminUser(ctx, a.Db)

	srv := &http.Server{
		Handler:      a.Router,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Listening on port %s", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func createInitialAdminUser(ctx context.Context, db *database.Db) {
	name := os.Getenv("ADMIN_NAME")
	email := os.Getenv("ADMIN_EMAIL")
	pw := os.Getenv("ADMIN_PASSWORD")

	// do not create admin if values are empty
	if util.ContainsEmpty(name, email, pw) {
		return
	}

	// check if admin already exists
	count, err := db.Users.CountDocuments(ctx, bson.D{{Key: "email", Value: email}})
	if err != nil {
		log.Fatal(err)
	}
	if count == 1 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	doc := &bson.D{
		{Key: "email", Value: email},
		{Key: "name", Value: name},
		{Key: "password", Value: string(hash)},
		{Key: "type", Value: models.AdminUser},
		{Key: "created_at", Value: time.Now()},
	}
	if _, err = db.Users.InsertOne(ctx, doc); err != nil {
		log.Fatal(err)
	}
}

// Function adds CORS middleware that allows cross origin requests when in development mode
func allowCrossOrigin(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Plaid-Public-Token")
		c.Next()
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
