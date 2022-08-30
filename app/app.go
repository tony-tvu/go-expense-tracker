package app

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
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/plaidapi"
	"github.com/tony-tvu/goexpense/user"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Init handlers
	u := &user.UserHandler{Db: a.Db}

	if env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.ForwardedByClientIP = true

	// apply global middleware
	router.Use(middleware.CORS(&env))
	router.Use(middleware.Logger(env))
	router.Use(middleware.RateLimit())
	router.Use(middleware.NoCache)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Ok"})
	})
	router.POST("/login", u.Login, middleware.LoginRateLimit())

	authRequired := router.Group("/", middleware.AuthRequired(a.Db))
	{
		authRequired.POST("/logout", u.Logout)
		authRequired.GET("/user_info", u.GetUserInfo)
		authRequired.GET("/create_link_token", plaidapi.CreateLinkToken)
		authRequired.POST("/set_access_token", plaidapi.SetAccessToken)

		adminRequired := authRequired.Group("/", middleware.AdminRequired(a.Db))
		{
			adminRequired.POST("/invite", u.InviteUser)
			adminRequired.GET("/sessions", u.GetSessions)
		}
	}

	// serve frontend
	router.Use(static.Serve("/", static.LocalFile("./web/build", true)))
	// prevent returning 404 when reloading page on frontend route
	router.NoRoute(func(ctx *gin.Context) {
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

	srv := &http.Server{
		Handler:      a.Router,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
