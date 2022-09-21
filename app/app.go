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
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/handlers"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/tasks"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Db           *database.MongoDb
	ConfigsCache *cache.Configs
	Router       *gin.Engine
	PlaidClient  *plaid.APIClient
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

	// Plaid
	var plaidEnvs = map[string]plaid.Environment{
		"sandbox":     plaid.Sandbox,
		"development": plaid.Development,
		"production":  plaid.Production,
	}
	clientID := os.Getenv("PLAID_CLIENT_ID")
	secret := os.Getenv("PLAID_SECRET")
	plaidEnv := os.Getenv("PLAID_ENV")
	if util.ContainsEmpty(clientID, secret, plaidEnv) {
		log.Println("plaid env configs are missing - service will not work")
	}
	plaidCfg := plaid.NewConfiguration()
	plaidCfg.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	plaidCfg.AddDefaultHeader("PLAID-SECRET", secret)
	plaidCfg.UseEnvironment(plaidEnvs[plaidEnv])
	pc := plaid.NewAPIClient(plaidCfg)
	a.PlaidClient = pc

	// Handlers
	users := &handlers.UserHandler{Db: a.Db}
	items := &handlers.ItemHandler{Db: a.Db, Client: pc}
	transactions := &handlers.TransactionHandler{Db: a.Db, ConfigsCache: a.ConfigsCache}
	configs := &handlers.ConfigsHandler{Db: a.Db, ConfigsCache: a.ConfigsCache}

	// Router
	if env == Production {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.ForwardedByClientIP = true
	if env == Development {
		allowCrossOrigin(router)
	}

	// middleware
	router.Use(middleware.RateLimit())
	router.Use(middleware.Logger(env))

	api := router.Group("/api", middleware.NoCache)
	{
		// configs
		api.GET("/registration_enabled", configs.RegistrationEnabled)
		api.GET("/configs", configs.GetConfigs)

		// users
		api.POST("/logout", users.Logout)
		api.POST("/login", middleware.LoginRateLimit(), users.Login)
		api.GET("/logged_in", users.IsLoggedIn)
		api.GET("/user_info", users.GetUserInfo)
		api.GET("/sessions", users.GetSessions)

		// items
		api.GET("/link_token", items.GetLinkToken)
		api.GET("/items", items.GetItems)
		api.POST("/items", items.CreateItem)
		api.DELETE("/items", items.DeleteItem)

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
	a.Db.Configs = mongoclient.Database(dbName).Collection("configs")
	a.Db.Items = mongoclient.Database(dbName).Collection("items")
	a.Db.Sessions = mongoclient.Database(dbName).Collection("sessions")
	a.Db.Transactions = mongoclient.Database(dbName).Collection("transactions")
	a.Db.Users = mongoclient.Database(dbName).Collection("users")

	// Create unique constraints
	if _, err = a.Db.Users.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Fatal(err)
	}
	if _, err = a.Db.Users.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Fatal(err)
	}
	if _, err = a.Db.Transactions.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "transaction_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Fatal(err)
	}
	createInitialAdminUser(ctx, a.Db)

	// Populate cache
	a.ConfigsCache.InitConfigsCache(ctx, a.Db)

	// Start scheduled tasks
	tasks.Start(a.Db, a.PlaidClient)

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

// Creates initial admin user. Account details can be specified in .env
func createInitialAdminUser(ctx context.Context, db *database.MongoDb) {
	username := os.Getenv("ADMIN_USERNAME")
	email := os.Getenv("ADMIN_EMAIL")
	pw := os.Getenv("ADMIN_PASSWORD")
	if util.ContainsEmpty(username, email, pw) {
		return
	}

	// check if admin already exists
	count, err := db.Users.CountDocuments(ctx, bson.D{{Key: "username", Value: username}})
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
		{Key: "username", Value: username},
		{Key: "email", Value: email},
		{Key: "password", Value: string(hash)},
		{Key: "type", Value: models.AdminUser},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}
	if _, err = db.Users.InsertOne(ctx, doc); err != nil {
		log.Fatal(err)
	}
}

// Allows cross origin requests from frontend server when in development
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
		MaxAge:           5 * time.Minute,
	}))
}
