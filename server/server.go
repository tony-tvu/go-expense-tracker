package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/plaid"
	"github.com/rs/cors"
	"github.com/tony-tvu/goexpense/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	App *app.App
}

func (s *Server) Init(ctx context.Context) {
	s.App = &app.App{}
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("failed to start - missing SECRET")
	}
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		log.Fatal("failed to start - missing JWT_KEY")
	}

	refreshTokenExpInt, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))
	if err != nil {
		refreshTokenExpInt = 86400
	}
	accessTokenExpInt, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	if err != nil {
		accessTokenExpInt = 900
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("failed to start - missing MONGODB_URI")
	}

	database := os.Getenv("DATABASE")
	if database == "" {
		database = "goexpense_local"
	}

	plaidClientID := os.Getenv("PLAID_CLIENT_ID")
	plaidSecret := os.Getenv("PLAID_SECRET")
	plaidEnv := os.Getenv("PLAID_ENV")
	plaidProducts := os.Getenv("PLAID_PRODUCTS")
	plaidCountryCodes := os.Getenv("PLAID_COUNTRY_CODES")
	if plaidClientID == "" || plaidSecret == "" || plaidEnv == "" || plaidProducts == "" || plaidCountryCodes == "" {
		log.Fatal("failed to start - missing Plaid env values")
	}
	plaidCfg := plaid.NewConfiguration()
	plaidCfg.AddDefaultHeader("PLAID-CLIENT-ID", plaidClientID)
	plaidCfg.AddDefaultHeader("PLAID-SECRET", plaidSecret)
	plaidCfg.UseEnvironment(app.PlaidEnvs[plaidEnv])
	plaidClient := plaid.NewAPIClient(plaidCfg)

	s.App = &app.App{
		Env:             env,
		Port:            port,
		Secret:          secret,
		JwtKey:          []byte(jwtKey),
		RefreshTokenExp: refreshTokenExpInt,
		AccessTokenExp:  accessTokenExpInt,
		MongoURI:        mongoURI,
		Db:              database,
		Coll: &app.Collections{
			Users:    "users",
			Sessions: "sessions",
		},
		PlaidClient:  plaidClient,
		CountryCodes: plaidCountryCodes,
		Products:     plaidProducts,
	}
	s.App.Router = InitRouter(ctx, s.App)
}

func (s *Server) Run(ctx context.Context) {
	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(s.App.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	s.App.MongoClient = mongoclient
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()

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

	log.Printf("Listening on port %s", s.App.Port)
	if err := http.ListenAndServe(":"+s.App.Port, handler); err != nil {
		log.Fatal(err)
	}
}
