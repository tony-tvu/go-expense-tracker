package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

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

func (s *Server) Init(ctx context.Context, env, port, authKey, jwtKey, refreshTokenExp, accessTokenExp, mongoURI, dbName, plaidClientID, plaidSecret, plaidEnv, plaidProducts, plaidountryCodes string) *mongo.Client {
	s.App = &app.App{}

	s.App.Env = env
	if env == "" {
		s.App.Env = "development"
	}
	s.App.Port = port
	if port == "" {
		s.App.Port = "8080"
	}
	s.App.AuthKey = authKey
	if authKey == "" {
		log.Fatal("failed to start - missing AUTH_KEY")
	}
	s.App.JwtKey = jwtKey
	if jwtKey == "" {
		log.Fatal("failed to start - missing JWT_KEY")
	}
	refreshTokenExpInt, err := strconv.Atoi(refreshTokenExp)
	s.App.RefreshTokenExp = refreshTokenExpInt
	if refreshTokenExp == "" || err != nil {
		s.App.RefreshTokenExp = 86400
	}
	accessTokenExpInt, err := strconv.Atoi(refreshTokenExp)
	s.App.AccessTokenExp = accessTokenExpInt
	if accessTokenExp == "" || err != nil {
		s.App.AccessTokenExp = 900
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
	router.HandleFunc("/api/health",
		Chain(handlers.Health, Middlewares...)).Methods("GET")
	router.Handle("/api/login_email",
		Chain(handlers.LoginEmail(ctx, s.App), Logging(), LoginRateLimit(), NoCache())).Methods("POST")
	router.Handle("/api/get_token_exp",
		Chain(handlers.IsTokenValid(ctx, s.App), Logging(), LoginRateLimit(), NoCache())).Methods("POST")
	router.Handle("/api/user",
		Chain(handlers.CreateUser(ctx, s.App), Middlewares...)).Methods("POST")
	router.Handle("/api/expense",
		Chain(handlers.GetExpenses(ctx, s.App), Middlewares...)).Methods("GET")
	// TODO: add auth to this so only registered users can create link tokens
	router.Handle("/api/create_link_token",
		Chain(handlers.CreateLinkToken(ctx, s.App), Middlewares...)).Methods("GET")
	router.Handle("/api/set_access_token",
		Chain(handlers.SetAccessToken(ctx, s.App), Middlewares...)).Methods("POST")
	router.PathPrefix("/").Handler(
		Chain(handlers.SpaHandler("web/build", "index.html"), Middlewares...)).Methods("GET")
	s.App.Router = router

	return mongoclient
}

func (s *Server) Start() {
	log.Printf("Listening on port %s", s.App.Port)

	var handler http.Handler
	if s.App.Env == "development" {
		handler = cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
			AllowedHeaders:   []string{"Content-Type", "Public-Token"},
			AllowCredentials: true,
		}).Handler(s.App.Router)
	} else {
		handler = s.App.Router
	}

	if err := http.ListenAndServe(":"+s.App.Port, handler); err != nil {
		log.Fatal(err)
	}
}
