package devtools

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Allows cross origin requests from frontend dev server when in development
func AllowCrossOrigin(r *gin.Engine) {
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

func CreateDummyUsers(ctx context.Context, db *gorm.DB) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	var user1Exists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?) AS found", "user1").Scan(&user1Exists)
	if !user1Exists {
		db.Create(&entity.User{
			Username: "user1",
			Email:    "user1@email.com",
			Password: string(hash),
			Type:     entity.RegularUser,
		})
	}

	var user2Exists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?) AS found", "user2").Scan(&user2Exists)
	if !user2Exists {
		db.Create(&entity.User{
			Username: "user2",
			Email:    "user2@email.com",
			Password: string(hash),
			Type:     entity.RegularUser,
		})
	}
}
