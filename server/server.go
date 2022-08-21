package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/plaid/plaid-go/plaid"
	"github.com/rs/cors"
	"github.com/tony-tvu/goexpense/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	App *app.App
}

func (s *Server) Init(
	ctx context.Context,
	port,
	env,
	secret,
	jwtKey,
	refreshTokenExp,
	accessTokenExp,
	mongoURI,
	database,
	plaidClientID,
	plaidSecret,
	plaidEnv,
	plaidCountryCodes,
	plaidProducts string,
) {
	s.App = &app.App{}

	if env == "" {
		env = "development"
	}
	if port == "" {
		port = "8080"
	}
	if secret == "" {
		log.Fatal("failed to start - missing SECRET")
	}
	if jwtKey == "" {
		log.Fatal("failed to start - missing JWT_KEY")
	}
	if mongoURI == "" {
		log.Fatal("failed to start - missing MONGODB_URI")
	}
	if database == "" {
		database = "goexpense_local"
	}

	refreshTokenExpInt, err := strconv.Atoi(refreshTokenExp)
	if err != nil {
		refreshTokenExpInt = 86400
	}
	accessTokenExpInt, err := strconv.Atoi(accessTokenExp)
	if err != nil {
		accessTokenExpInt = 900
	}

	if plaidClientID == "" || plaidSecret == "" || plaidEnv == "" || plaidProducts == "" || plaidCountryCodes == "" {
		log.Println("plaid configs are missing - service will not work")
	}
	plaidCfg := plaid.NewConfiguration()
	plaidCfg.AddDefaultHeader("PLAID-CLIENT-ID", plaidClientID)
	plaidCfg.AddDefaultHeader("PLAID-SECRET", plaidSecret)
	plaidCfg.UseEnvironment(app.PlaidEnvs[plaidEnv])
	plaidClient := plaid.NewAPIClient(plaidCfg)
	s.App.PlaidClient = plaidClient

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
		PlaidClientID:     plaidClientID,
		PlaidSecret:       plaidSecret,
		PlaidEnv:          plaidEnv,
		PlaidCountryCodes: plaidCountryCodes,
		PlaidProducts:     plaidProducts,
	}
	s.App.Router = InitRouter(ctx, s.App)

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

	var handler http.Handler
	if s.App.Env == "development" {
		handler = cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
			AllowedHeaders:   []string{"Content-Type", "Public-Token"},
			AllowCredentials: true,
		}).Handler(s.App.Router)
	} else {
		handler = s.App.Router
	}

	srv := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf(":%s", s.App.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}
	s.App.Server = srv
}
