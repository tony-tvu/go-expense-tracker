package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/plaid/plaid-go/plaid"
	"github.com/rs/cors"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	App *app.App
}

func (s *Server) Initialize(ctx context.Context) {
	if s.App.Env == "" {
		s.App.Env = "development"
	}
	if s.App.Port == "" {
		s.App.Port = "8080"
	}
	if s.App.Secret == "" {
		log.Fatal("failed to start - missing SECRET")
	}
	if string(s.App.JwtKey) == "" {
		log.Fatal("failed to start - missing JWT_KEY")
	}
	if s.App.MongoURI == "" {
		log.Fatal("failed to start - missing MONGODB_URI")
	}
	if s.App.Db == "" {
		s.App.Db = "goexpense_local"
	}

	// Plaid client
	if s.App.PlaidClientID != "" ||
		s.App.PlaidSecret != "" ||
		s.App.PlaidEnv != "" ||
		s.App.PlaidProducts != "" ||
		s.App.PlaidCountryCodes != "" {
		plaidCfg := plaid.NewConfiguration()
		plaidCfg.AddDefaultHeader("PLAID-CLIENT-ID", s.App.PlaidClientID)
		plaidCfg.AddDefaultHeader("PLAID-SECRET", s.App.PlaidSecret)
		plaidCfg.UseEnvironment(app.PlaidEnvs[s.App.PlaidEnv])
		plaidClient := plaid.NewAPIClient(plaidCfg)
		s.App.PlaidClient = plaidClient
	} else {
		log.Println("plaid configs are missing - service will not work")
	}

	// Mongo collections
	s.App.Coll = &app.Collections{
		Users:    "users",
		Sessions: "sessions",
	}

	// Handlers
	handlers := &app.Handlers{
		// Finances
		GetExpenses: Chain(handler.GetExpenses(s.App), UseMiddlewares()...),
		// Health
		Health: Chain(handler.Health, UseMiddlewares()...),
		// Plaid
		// TODO: add auth to this so only registered users can create link tokens
		CreateLinkToken: Chain(handler.CreateLinkToken(s.App), UseMiddlewares()...),
		SetAccessToken:  Chain(handler.SetAccessToken(s.App), UseMiddlewares()...),
		// Users
		CreateUser:  Chain(handler.CreateUser(s.App), UseMiddlewares()...),
		GetUserInfo: Chain(handler.GetUserInfo(s.App), UseMiddlewares()...),
		LoginEmail:  Chain(handler.LoginEmail(s.App), UseMiddlewares(LoginRateLimit())...),
		// Web
		ServeClient: Chain(handler.ServeClient("web/build", "index.html"), UseMiddlewares()...),
	}
	s.App.Handlers = handlers

	// Routes
	router := mux.NewRouter()
	// Finances
	router.Handle("/api/expense", s.App.Handlers.GetExpenses).Methods("GET")
	// Health
	router.HandleFunc("/api/health", s.App.Handlers.Health).Methods("GET")
	// Users
	router.Handle("/api/create_user", s.App.Handlers.CreateUser).Methods("POST")
	router.Handle("/api/login_email", s.App.Handlers.LoginEmail).Methods("POST")
	router.Handle("/api/user_info", s.App.Handlers.GetUserInfo).Methods("GET")
	// Plaid
	router.Handle("/api/create_link_token", s.App.Handlers.CreateLinkToken).Methods("GET")
	router.Handle("/api/set_access_token", s.App.Handlers.CreateLinkToken).Methods("POST")
	// Web
	router.PathPrefix("/").Handler(s.App.Handlers.CreateUser).Methods("GET")

	var h http.Handler
	if s.App.Env == "development" {
		h = cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
			AllowedHeaders:   []string{"Content-Type", "Public-Token"},
			AllowCredentials: true,
		}).Handler(router)
	} else {
		h = router
	}

	srv := &http.Server{
		Handler:      h,
		Addr:         fmt.Sprintf(":%s", s.App.Port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	s.App.Server = srv
}

func (s *Server) Run(ctx context.Context) {
	if s.App.Env != "test" {
		mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(s.App.MongoURI))
		if err != nil {
			log.Fatal(err)
		}
		err = mongoclient.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := mongoclient.Disconnect(ctx); err != nil {
				log.Println("mongo has been disconnected: ", err)
			}
		}()
		s.App.MongoClient = mongoclient
	}

	log.Printf("Listening on port %s", s.App.Port)
	err := s.App.Server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
