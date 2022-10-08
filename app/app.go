package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/handlers"
	"github.com/tony-tvu/goexpense/jobs"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Db           *database.MongoDb
	ConfigsCache *cache.Configs
	Router       *gin.Engine
	Jobs         *jobs.Jobs
}

const (
	Production  = "production"
	Development = "development"
)

var env string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}
}

func (a *App) Initialize(ctx context.Context) {
	env = os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV is not set")
	}
	a.Db = &database.MongoDb{}
	a.ConfigsCache = &cache.Configs{}

	// Jobs
	jobs := &jobs.Jobs{Db: a.Db}
	jobsEnabled, err := strconv.ParseBool(os.Getenv("JOBS_ENABLED"))
	if err != nil {
		jobs.Enabled = false
	} else {
		jobs.Enabled = jobsEnabled
	}
	jobsInterval, err := strconv.Atoi(os.Getenv("JOBS_INTERVAL"))
	if err != nil {
		jobs.Interval = 86400
	} else {
		jobs.Interval = jobsInterval
	}
	a.Jobs = jobs

	// TellerClient
	dirname, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	certPath := path.Join(dirname, "/certificate/certificate.pem")
	keyPath := path.Join(dirname, "/certificate/private_key.pem")
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}
	tellerClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
		Timeout: 2 * time.Minute,
	}

	// Handlers
	users := &handlers.UserHandler{Db: a.Db}
	teller := &handlers.TellerHandler{Db: a.Db, ConfigsCache: a.ConfigsCache, TellerClient: tellerClient}
	configs := &handlers.ConfigsHandler{Db: a.Db, ConfigsCache: a.ConfigsCache}

	// Router
	if env == Production {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.ForwardedByClientIP = true

	// Cors
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	allowedOriginsArr := strings.Split(allowedOrigins, ",")
	router.Use(middleware.CorsHeaders(allowedOriginsArr))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOriginsArr,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
		MaxAge:           5 * time.Minute,
	}))

	// middleware
	router.Use(middleware.RateLimit())
	router.Use(middleware.Logger(env))

	api := router.Group("/api", middleware.NoCache)
	{
		// configs
		api.GET("/registration_enabled", configs.RegistrationEnabled)
		api.GET("/teller_app_id", configs.TellerAppID)
		api.GET("/configs", configs.GetConfigs)
		api.PUT("/configs", configs.UpdateConfigs)

		// users
		api.POST("/logout", users.Logout)
		api.POST("/login", middleware.LoginRateLimit(), users.Login)
		api.GET("/logged_in", users.IsLoggedIn)
		api.GET("/user_info", users.GetUserInfo)
		api.GET("/sessions", users.GetSessions)

		// teller
		api.POST("/enrollments", teller.NewEnrollment)
		api.GET("/enrollments", teller.GetEnrollments)
	}

	router.Use(middleware.FrontendCache, static.Serve("/", static.LocalFile("./web/build", true)))
	router.NoRoute(middleware.FrontendCache, func(ctx *gin.Context) {
		ctx.File("./web/build")
	})
	a.Router = router
}

func (a *App) Start(ctx context.Context) {
	// Start mongodb
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	if util.ContainsEmpty(mongoURI, dbName) {
		log.Fatal("env variables are missing")
	}
	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()
	a.Db.SetCollections(mongoclient, dbName)
	a.Db.CreateUniqueConstraints(ctx)
	username := os.Getenv("ADMIN_USERNAME")
	email := os.Getenv("ADMIN_EMAIL")
	pw := os.Getenv("ADMIN_PASSWORD")
	if util.ContainsEmpty(username, email, pw) {
		return
	}
	a.Db.CreateInitialAdminUser(ctx, username, email, pw)

	// Populate cache
	a.ConfigsCache.InitConfigsCache(ctx, a.Db)

	// Start scheduled jobs
	a.Jobs.Start(ctx)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	log.Printf("Listening on port %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
