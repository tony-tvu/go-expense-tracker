package app

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tony-tvu/goexpense/config"
	"github.com/tony-tvu/goexpense/routes"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Config      config.AppConfig
	MongoClient *mongo.Client
	Router      *mux.Router
}

func (a *App) Initialize(ctx context.Context, env, port, authKey, mongoURI, dbName string) *mongo.Client {
	cfg := &config.AppConfig{}
	cfg.Env = env
	if env == "" {
		cfg.Env = "DEV"
	}
	cfg.Port = port
	if port == "" {
		cfg.Port = "8080"
	}
	cfg.AuthKey = authKey
	if authKey == "" {
		log.Fatal("failed to start - missing AUTH_KEY")
	}
	cfg.MongoURI = mongoURI
	if authKey == "" {
		log.Fatal("failed to start - missing MONGODB_URI")
	}
	cfg.DbName = dbName
	if dbName == "" {
		cfg.DbName = "goexpense"
	}

	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	a.MongoClient = mongoclient
	cfg.UserCollection = "users"

	a.Config = *cfg
	a.Router = routes.GetRouter(cfg, a.MongoClient)

	return mongoclient
}

func (a *App) Run() {
	log.Printf("Listening on port %s", a.Config.Port)
	if err := http.ListenAndServe(":"+a.Config.Port, a.Router); err != nil {
		log.Fatal(err)
	}
}
