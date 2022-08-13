package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/config"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/user"
	"github.com/tony-tvu/goexpense/web"

	"github.com/gorilla/mux"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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
	dbName := os.Getenv("DATABASE_NAME")
	dbTimeout, _ := strconv.Atoi(os.Getenv("DB_TIMEOUT_SECONDS"))
	authKeyStr := os.Getenv("KEY")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// TODO: throw error if these are blank

	// Setup MongoDB
	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Setup Rate Limiter
	rate, err := limiter.NewRateFromFormatted("3000-M")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	rateLimit := stdlib.NewMiddleware(instance)

	// Init Config
	cfg := config.Config{
		Client:         mongoclient,
		Database:       dbName,
		DBTimeout:      dbTimeout,
		UserCollection: "users",
		AuthKey:        []byte(authKeyStr),
	}

	// Handlers
	uh := user.Handler{Config: cfg}

	// Routes
	router := mux.NewRouter()

	router.HandleFunc("/api/health", middleware.Chain(HealthHandler, middleware.CommonMiddleware...))

	router.Handle("/api/user", rateLimit.Handler(http.HandlerFunc(uh.NewHandler)))
	router.Handle("/", rateLimit.Handler(web.SpaHandler{StaticPath: "web/build", IndexPath: "index.html"}))

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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		body := make(map[string]string)
		body["message"] = "Ok"
		jData, _ := json.Marshal(body)
		w.Write(jData)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
