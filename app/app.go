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
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/handlers"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/tasks"
	"github.com/tony-tvu/goexpense/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Db          *gorm.DB
	Router      *gin.Engine
	PlaidClient *plaid.APIClient
}

const (
	Production  = "production"
	Development = "development"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}
}

func (a *App) Initialize(ctx context.Context) {
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV is not set")
	}

	// Database
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	if util.ContainsEmpty(dbUser, dbPwd, dbHost, dbName) {
		log.Fatal("postgres config envs are missing")
	}

	var dbURI string
	if dbPort == "" {
		dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s",
			dbUser, dbPwd, dbName, dbHost)
	} else {
		dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s port=%s",
			dbUser, dbPwd, dbName, dbHost, dbPort)
	}

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&entity.Session{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Item{})
	db.AutoMigrate(&entity.Transaction{})
	db.AutoMigrate(&entity.Config{})

	createInitialAdminUser(ctx, db)
	a.Db = db

	// Caches
	cachedConf := &cache.Configs{}
	cachedConf.InitConfigsCache(db)

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
	users := &handlers.UserHandler{Db: db}
	items := &handlers.ItemHandler{Db: db, Client: pc}
	transactions := &handlers.TransactionHandler{Db: db}
	configs := &handlers.ConfigsHandler{Cache: cachedConf}

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
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	})

	api := router.Group("/api", middleware.NoCache)
	{
		// configs
		api.GET("/registration_enabled", configs.RegistrationEnabled)

		// users
		api.POST("/logout", users.Logout)
		api.POST("/login", middleware.LoginRateLimit(), users.Login)
		api.GET("/logged_in", users.IsLoggedIn)
		api.GET("/is_admin", users.IsAdmin)
		api.GET("/user_info", users.GetUserInfo)
		api.GET("/sessions", users.GetSessions)

		// items
		api.GET("/link_token", items.GetLinkToken)
		api.GET("/items", items.GetItems)
		api.POST("/items", items.CreateItem)
		api.DELETE("/items", items.DeleteItem)

		// transactions
		api.GET("/transactions", transactions.GetTransactions)
	}

	router.Use(middleware.FrontendCache, static.Serve("/", static.LocalFile("./web/build", true)))
	router.NoRoute(middleware.FrontendCache, func(ctx *gin.Context) {
		ctx.File("./web/build")
	})
	a.Router = router
}

func (a *App) Serve() {
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

func createInitialAdminUser(ctx context.Context, db *gorm.DB) {
	username := os.Getenv("ADMIN_USERNAME")
	email := os.Getenv("ADMIN_EMAIL")
	pw := os.Getenv("ADMIN_PASSWORD")
	if util.ContainsEmpty(username, email, pw) {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	var exists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?) AS found", username).Scan(&exists)
	if !exists {
		db.Create(&entity.User{
			Username: username,
			Email:    email,
			Password: string(hash),
			Type:     entity.AdminUser,
		})
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
