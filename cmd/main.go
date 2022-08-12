package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/web"

	"github.com/tony-tvu/goexpense/user"

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
	addr := os.Getenv("ADDR")
	dbName := os.Getenv("DATABASE_NAME")
	authKeyStr := os.Getenv("KEY")

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
	user := &user.UserConfigs{
		Client:     mongoclient,
		Database:   dbName,
		Collection: "users",
		AuthKey: []byte(authKeyStr),
	}

	// Setup Rate Limiter
	rate, err := limiter.NewRateFromFormatted("3000-M")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	rateLimit := stdlib.NewMiddleware(instance)

	// Routes
	router := mux.NewRouter()
	router.Handle("/api/health", rateLimit.Handler(http.HandlerFunc(HealthHandler)))
	router.Handle("/api/user", rateLimit.Handler(http.HandlerFunc(user.Handler)))
	router.Handle("/", rateLimit.Handler(web.SpaHandler{StaticPath: "web/build", IndexPath: "index.html"}))

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server started on %s\n", addr)
	log.Fatal(srv.ListenAndServe())
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprint(w, "Ok")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
