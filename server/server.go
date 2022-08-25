package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/plaid/plaid-go/plaid"
	"github.com/rs/cors"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/plaidapi"
	"github.com/tony-tvu/goexpense/user"
	"github.com/tony-tvu/goexpense/web"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	App *app.App
}

func (s *Server) Initialize() {
	if string(s.App.JwtKey) == "" {
		log.Fatal("failed to start - missing JWT_KEY")
	}
	if s.App.MongoURI == "" {
		log.Fatal("failed to start - missing MONGODB_URI")
	}
	if s.App.Env == "" {
		s.App.Env = "development"
	}
	if s.App.Port == "" {
		s.App.Port = "8080"
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

	// Handlers
	handlers := &app.Handlers{
		// Auth
		EmailLogin: Chain(auth.EmailLogin(s.App), UseMiddlewares(s.App, LoginRateLimit())...),
		// Health
		Health: Chain(Health, UseMiddlewares(s.App)...),
		// Plaid
		// TODO: add auth to this so only registered users can create link tokens
		CreateLinkToken: Chain(plaidapi.CreateLinkToken(s.App), UseMiddlewares(s.App)...),
		SetAccessToken:  Chain(plaidapi.SetAccessToken(s.App), UseMiddlewares(s.App)...),
		// Users
		CreateUser:  Chain(user.Create(s.App), UseMiddlewares(s.App)...),
		GetUserInfo: Chain(user.GetInfo(s.App), UseMiddlewares(s.App, LoginProtected(s.App))...),
		// Web
		ServeClient: Chain(web.Serve("web/build", "index.html"), UseMiddlewares(s.App)...),
	}
	s.App.Handlers = handlers

	// Routes
	router := mux.NewRouter()
	// Auth
	router.Handle("/api/email_login", s.App.Handlers.EmailLogin).Methods("POST")
	// Finances
	router.Handle("/api/expense", s.App.Handlers.GetExpenses).Methods("GET")
	// Health
	router.HandleFunc("/api/health", s.App.Handlers.Health).Methods("GET")
	// Users
	router.Handle("/api/create_user", s.App.Handlers.CreateUser).Methods("POST")
	router.Handle("/api/user_info", s.App.Handlers.GetUserInfo).Methods("GET")
	// Plaid
	router.Handle("/api/create_link_token", s.App.Handlers.CreateLinkToken).Methods("GET")
	router.Handle("/api/set_access_token", s.App.Handlers.CreateLinkToken).Methods("POST")
	// Web
	router.PathPrefix("/").Handler(s.App.Handlers.CreateUser).Methods("GET")

	s.App.Router = router
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
		if s.App.DbName == "" {
			s.App.DbName = "goexpense_local"
		}
		s.App.Users = mongoclient.Database(s.App.DbName).Collection("users")
		s.App.Users = mongoclient.Database(s.App.DbName).Collection("sessions")
	}

	var h http.Handler
	if s.App.Env == "development" {
		h = cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
			AllowedHeaders:   []string{"Content-Type", "Plaid-Public-Token", "Google-ID-Token"},
			AllowCredentials: true,
		}).Handler(s.App.Router)
	} else {
		h = cors.New(cors.Options{
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
			AllowedHeaders:   []string{"Content-Type", "Plaid-Public-Token", "Google-ID-Token"},
			AllowCredentials: true,
		}).Handler(s.App.Router)
	}

	srv := &http.Server{
		Handler:      h,
		Addr:         fmt.Sprintf(":%s", s.App.Port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Printf("Listening on port %s", s.App.Port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := make(map[string]string)
	body["message"] = "Ok"
	jData, _ := json.Marshal(body)
	w.Write(jData)
}
