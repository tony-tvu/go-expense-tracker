package app

import (
	"github.com/gorilla/mux"
	"github.com/plaid/plaid-go/plaid"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Env               string
	Port              string
	EncryptionKey     string
	JwtKey            string
	RefreshTokenExp   int
	AccessTokenExp    int
	MongoURI          string
	DbName            string
	Users             *mongo.Collection
	Sessions          *mongo.Collection
	PlaidClient       *plaid.APIClient
	PlaidClientID     string
	PlaidSecret       string
	PlaidEnv          string
	PlaidCountryCodes string
	PlaidProducts     string
	Router            *mux.Router
}

var PlaidEnvs = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}
