package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	if string(s.App.EncryptionKey) == "" {
		log.Fatal("fatal: missing ENCRYPTION_KEY")
	}
	if string(s.App.JwtKey) == "" {
		log.Fatal("fatal: missing JWT_KEY")
	}
	if s.App.MongoURI == "" {
		log.Fatal("fatal: missing MONGODB_URI")
	}
	if s.App.Env == "" {
		s.App.Env = "development"
	}
	if s.App.Port == "" {
		s.App.Port = "80"
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
	authHandler := &auth.AuthHandler{App: s.App}
	plaidHandler := &plaidapi.PlaidHandler{App: s.App}
	userHandler := &user.UserHandler{App: s.App}
	spaHandler := &web.SpaHandler{StaticPath: "web/build", IndexPath: "index.html"}

	// Routes
	router := mux.NewRouter()

	// Auth
	router.Handle("/api/login", Chain(authHandler.Login,
		Middlewares(s.App, Method("POST"), LoginRateLimit())...))
	router.Handle("/api/logout", Chain(authHandler.Logout,
		Middlewares(s.App, Method("DELETE"))...))
	router.Handle("/api/sessions", Chain(authHandler.GetSessions,
		Middlewares(s.App, Method("GET"), LoggedIn(s.App), Admin(s.App))...))

	// Health
	router.Handle("/api/health", Chain(Health,
		Middlewares(s.App, Method("GET"))...))

	// Users
	router.Handle("/api/user_info", Chain(userHandler.GetInfo,
		Middlewares(s.App, Method("GET"), LoggedIn(s.App))...))

	// If true, the '/api/create_user' endpoint is available to everyone.
	// If false, only admins can create new users.
	isAllowed, err := strconv.ParseBool(s.App.AllowCreateUsers)
	if err != nil {
		isAllowed = false
	}
	if isAllowed {
		router.Handle("/api/create_user", Chain(userHandler.Create,
			Middlewares(s.App, Method("POST"))...))
	} else {
		router.Handle("/api/create_user", Chain(userHandler.Create,
			Middlewares(s.App, Method("POST"), LoggedIn(s.App), Admin(s.App))...))
	}

	// Plaid
	router.Handle("/api/create_link_token", Chain(plaidHandler.CreateLinkToken,
		Middlewares(s.App, Method("GET"), LoggedIn(s.App))...))
	router.Handle("/api/set_access_token", Chain(plaidHandler.SetAccessToken,
		Middlewares(s.App, Method("POST"), LoggedIn(s.App))...))

	// Web
	router.PathPrefix("/").Handler(Chain(spaHandler.Serve,
		Middlewares(s.App, Method("GET"))...))

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
		s.App.Sessions = mongoclient.Database(s.App.DbName).Collection("sessions")
	}

	var h http.Handler
	if s.App.Env == "development" {
		h = cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"*"},
			AllowedHeaders:   []string{"Content-Type", "Plaid-Public-Token"},
			AllowCredentials: true,
		}).Handler(s.App.Router)
	} else {
		h = cors.New(cors.Options{
			AllowedMethods:   []string{"*"},
			AllowedHeaders:   []string{"Content-Type", "Plaid-Public-Token"},
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
