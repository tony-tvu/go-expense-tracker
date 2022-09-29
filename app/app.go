package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/tasks"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Db           *database.MongoDb
	ConfigsCache *cache.Configs
	Router       *gin.Engine
	Tasks        *tasks.Tasks
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

	// Tasks
	tasks := &tasks.Tasks{Db: a.Db}
	taskInterval, err := strconv.Atoi(os.Getenv("TASK_INTERVAL"))
	if err != nil {
		tasks.TaskInterval = 86400
	} else {
		tasks.TaskInterval = taskInterval
	}
	a.Tasks = tasks

	// Handlers
	users := &handlers.UserHandler{Db: a.Db}
	items := &handlers.ItemHandler{Db: a.Db, ConfigsCache: a.ConfigsCache, Tasks: tasks, WebhooksURL: os.Getenv("WEBHOOKS_URL")}
	transactions := &handlers.TransactionHandler{Db: a.Db, ConfigsCache: a.ConfigsCache}
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
		api.GET("/configs", configs.GetConfigs)
		api.PUT("/configs", configs.UpdateConfigs)

		// users
		api.POST("/logout", users.Logout)
		api.POST("/login", middleware.LoginRateLimit(), users.Login)
		api.GET("/logged_in", users.IsLoggedIn)
		api.GET("/user_info", users.GetUserInfo)
		api.GET("/sessions", users.GetSessions)

		// items
		api.GET("/link_token", items.GetLinkToken)
		api.GET("/items/:page", items.GetItems)
		api.POST("/items", items.CreateItem)
		api.DELETE("/items/:plaid_item_id", items.DeleteItem)
		api.GET("/accounts", items.GetAccounts)
		api.POST("/receive_webhooks", items.ReceiveWebooks)
		api.PUT("/update_webhooks_url", items.UpdateItemsWebhooksURL)

		// transactions
		api.GET("/transactions/:page", transactions.GetTransactions)
	}

	router.Use(middleware.FrontendCache, static.Serve("/", static.LocalFile("./web/build", true)))
	router.NoRoute(middleware.FrontendCache, func(ctx *gin.Context) {
		ctx.File("./web/build")
	})
	a.Router = router
}

func (a *App) Start(ctx context.Context) {
	// Start mongodb
	mongoURI := os.Getenv("MONGODB_URI")
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
	a.Db.Accounts = mongoclient.Database(dbName).Collection("accounts")
	a.Db.Configs = mongoclient.Database(dbName).Collection("configs")
	a.Db.Items = mongoclient.Database(dbName).Collection("items")
	a.Db.Sessions = mongoclient.Database(dbName).Collection("sessions")
	a.Db.Transactions = mongoclient.Database(dbName).Collection("transactions")
	a.Db.Users = mongoclient.Database(dbName).Collection("users")

	database.CreateUniqueConstraints(ctx, a.Db)
	username := os.Getenv("ADMIN_USERNAME")
	email := os.Getenv("ADMIN_EMAIL")
	pw := os.Getenv("ADMIN_PASSWORD")
	if util.ContainsEmpty(username, email, pw) {
		return
	}
	database.CreateInitialAdminUser(ctx, a.Db, username, email, pw)

	// Populate cache
	a.ConfigsCache.InitConfigsCache(ctx, a.Db)

	// Start scheduled tasks
	a.Tasks.Start(ctx)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on port %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
