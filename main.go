package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/handlers"
)

func main() {
	ctx := context.Background()

	// Get environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	dbName := os.Getenv("DATABASE_NAME")
	authKeyStr := os.Getenv("KEY")
	env := os.Getenv("ENV")
	if dbName == "" || authKeyStr == "" {
		log.Fatal("Missing environment variables")
	}

	// Create AppConfigs
	cfg := app.AppConfigs{
		Env:            env,
		Database:       dbName,
		DBTimeout:      5,
		UserCollection: "users",
		AuthKey:        []byte(authKeyStr),
	}

	// Create mongo client
	mc := database.GetMongoClient(ctx)
	defer func() {
		if err := mc.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Routes
	router := mux.NewRouter()

	router.HandleFunc("/api/health", app.Chain(handlers.HealthHandler, app.Middlewares...))
	router.Handle("/api/user", app.Chain(handlers.UserHandler(cfg, mc), app.Middlewares...))
	router.PathPrefix("/").Handler(app.Chain(handlers.SpaHandler("web/build", "index.html"), app.Middlewares...))

	port := "8080"
	srv := &http.Server{
		Handler:           router,
		Addr:              fmt.Sprintf(":%s", port),
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
