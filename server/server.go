package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/plaid/plaid-go/plaid"
	"github.com/rs/cors"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/handlers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	App *app.App
}

func (s *Server) Init(ctx context.Context, env, port, authKey, mongoURI, dbName, plaidClientID, plaidSecret, plaidEnv, plaidProducts, plaidountryCodes string) *mongo.Client {
	s.App = &app.App{}

	s.App.Env = env
	if env == "" {
		s.App.Env = "DEV"
	}
	s.App.Port = port
	if port == "" {
		s.App.Port = "8080"
	}
	s.App.AuthKey = authKey
	if authKey == "" {
		log.Fatal("failed to start - missing AUTH_KEY")
	}
	s.App.MongoURI = mongoURI
	if mongoURI == "" {
		log.Fatal("failed to start - missing MONGODB_URI")
	}
	s.App.DbName = dbName
	if dbName == "" {
		s.App.DbName = "goexpense"
	}
	if plaidClientID == "" || plaidSecret == "" || plaidEnv == "" || plaidProducts == "" || plaidountryCodes == "" {
		log.Fatal("failed to start - missing Plaid env values")
	}
	plaidCfg := plaid.NewConfiguration()
	plaidCfg.AddDefaultHeader("PLAID-CLIENT-ID", plaidClientID)
	plaidCfg.AddDefaultHeader("PLAID-SECRET", plaidSecret)
	plaidCfg.UseEnvironment(app.PlaidEnvs[plaidEnv])
	plaidApiClient := plaid.NewAPIClient(plaidCfg)
	pc := &app.PlaidClient{
		ClientID:     plaidClientID,
		Secret:       plaidSecret,
		Env:          plaidEnv,
		Products:     plaidProducts,
		CountryCodes: plaidountryCodes,
		ApiClient:    plaidApiClient,
	}
	s.App.PlaidClient = pc

	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(s.App.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	s.App.MongoClient = mongoclient
	s.App.UserCollection = "users"

	router := mux.NewRouter()
	router.HandleFunc("/api/health", Chain(handlers.Health, Middlewares...)).Methods("GET")
	router.Handle("/api/user", Chain(handlers.CreateUser(s.App), Middlewares...)).Methods("POST")
	router.Handle("/api/expense", Chain(handlers.GetExpenses(s.App), Middlewares...)).Methods("GET")
	// TODO: add auth to this so only registered users can create link tokens
	router.Handle("/api/create_link_token", Chain(handlers.CreateLinkToken(s.App), Middlewares...)).Methods("GET")
	router.PathPrefix("/").Handler(Chain(handlers.SpaHandler("web/build", "index.html"), Middlewares...)).Methods("GET")
	s.App.Router = router

	return mongoclient
}

func (s *Server) Start() {
	log.Printf("Listening on port %s", s.App.Port)

	var handler http.Handler
	if s.App.Env == "DEV" {
		handler = cors.Default().Handler(s.App.Router)
	} else {
		handler = s.App.Router
	}

	if err := http.ListenAndServe(":"+s.App.Port, handler); err != nil {
		log.Fatal(err)
	}
}
