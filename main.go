package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// "github.com/gin-gonic/contrib/static"
	// "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "github.com/tony-tvu/goexpense/user"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	// Get environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	// dbName := os.Getenv("DATABASE_NAME")
	// port := os.Getenv("PORT")

	// Init mongo
	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// userService := &user.UserService{
	// 	Client:     mongoclient,
	// 	Database:   dbName,
	// 	Collection: "users",
	// }

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("dist")))

	srv := &http.Server{
		Handler:      mux,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	
	// serve
	fmt.Println("Server started on PORT 8080")
	log.Fatal(srv.ListenAndServe())

	// router := gin.Default()

	// // Serve frontend
	// router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	// // Endpoints
	// api := router.Group("/api")
	// {
	// 	api.GET("/", func(c *gin.Context) {
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"message": "goexpense api",
	// 		})
	// 	})

	// 	api.POST("/user", func(c *gin.Context) {
	// 		c.String(http.StatusOK, userService.CreateUser(c))
	// 	})
	// }

	// router.Run(fmt.Sprint(":", port))
}
