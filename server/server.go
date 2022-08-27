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
	ah := &auth.AuthHandler{App: s.App}
	ph := &plaidapi.PlaidHandler{App: s.App}
	uh := &user.UserHandler{App: s.App}
	sh := &web.SpaHandler{StaticPath: "web/build", IndexPath: "index.html"}

	// Routes
	r := mux.NewRouter()
	// Auth
	r.Handle("/api/login", GuestUser(ah.Login, s.App)).Methods("POST")
	r.Handle("/api/logout", RegularUser(ah.Logout, s.App)).Methods("POST")
	r.Handle("/api/sessions", AdminUser(ah.GetSessions, s.App)).Methods("GET")
	// Health
	r.Handle("/api/health", GuestUser(Health, s.App)).Methods("GET")
	// Users
	r.Handle("/api/user_info", RegularUser(uh.GetInfo, s.App)).Methods("GET")
	r.Handle("/api/invite", AdminUser(uh.Invite, s.App)).Methods("POST")
	// Plaid
	r.Handle("/api/create_link_token", RegularUser(ph.CreateLinkToken, s.App)).Methods("GET")
	r.Handle("/api/set_access_token", RegularUser(ph.SetAccessToken, s.App)).Methods("POST")
	// Web
	r.PathPrefix("/").Handler(GuestUser(sh.Serve, s.App)).Methods("GET")

	s.App.Router = r
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

	h := cors.New(cors.Options{
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"Content-Type", "Plaid-Public-Token"},
		AllowCredentials: true,
	}).Handler(s.App.Router)

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
