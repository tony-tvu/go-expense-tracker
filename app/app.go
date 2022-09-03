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
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/plaidapi"
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

var env string
var port string
var dbURL string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}
	env = os.Getenv("ENV")
	port = os.Getenv("PORT")
	dbURL = os.Getenv("DB_URL")
	if util.ContainsEmpty(env, port, dbURL) {
		log.Fatal("env variables are missing")
	}
}

func (a *App) Initialize(ctx context.Context) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&entity.Session{})
	db.AutoMigrate(&entity.User{})
	a.Db = db

	u := &user.UserHandler{Db: db}

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
	createInitialAdminUser(ctx, a.Db)

	srv := &http.Server{
		Handler:      a.Router,
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

func createInitialAdminUser(ctx context.Context, db *gorm.DB) {
	name := os.Getenv("ADMIN_NAME")
	email := os.Getenv("ADMIN_EMAIL")
	pw := os.Getenv("ADMIN_PASSWORD")

	// do not create admin if values are empty
	if util.ContainsEmpty(name, email, pw) {
		return
	}

	// check if admin already exists
	var u *entity.User
	if result := db.Where("email = ?", email).First(&u); result.Error != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		if result := db.Create(&entity.User{
			Name:     name,
			Email:    email,
			Password: string(hash),
			Type:     entity.AdminUser,
		}); result.Error != nil {
			log.Fatal(err)
		}
	}

	// var res struct {
	// 	Found bool
	// }
	// db.Raw("SELECT EXISTS(SELECT 1 FROM user WHERE email = ?) AS found", email).Scan(&res)
	// if res.Found {
	// 	return
	// }
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
