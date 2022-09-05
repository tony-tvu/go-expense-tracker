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
	"github.com/tony-tvu/goexpense/devtools"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/plaidapi"
	"github.com/tony-tvu/goexpense/tasks"
	"github.com/tony-tvu/goexpense/user"
	"github.com/tony-tvu/goexpense/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Db     *gorm.DB
	Router *gin.Engine
}

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

	// Init database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("postgres url is missing")
	}
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&entity.Session{})
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.Item{})
	db.AutoMigrate(&entity.Transaction{})

	createInitialAdminUser(ctx, db)
	a.Db = db

	// Init handlers
	u := &user.UserHandler{Db: db}
	p := &plaidapi.PlaidHandler{Db: db}

	// Init router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.ForwardedByClientIP = true
	if env == "development" {
		devtools.AllowCrossOrigin(router)
		devtools.CreateDummyUsers(ctx, db)
	}

	// Apply middlewares
	router.Use(middleware.RateLimit())
	router.Use(middleware.Logger(env))

	// Declare routes
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
			authGroup.GET("/create_link_token", p.CreateLinkToken)
			authGroup.POST("/set_access_token", p.SetAccessToken)

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

func (a *App) Serve() {
	tasks.Start(a.Db)

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

	// do not create admin if values are empty
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
