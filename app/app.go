package app

import (
	"github.com/gorilla/mux"
	"github.com/plaid/plaid-go/plaid"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Env            string
	Port           string
	AuthKey        string
	MongoURI       string
	DbName         string
	UserCollection string
	MongoClient    *mongo.Client
	PlaidClient    *PlaidClient
	Router         *mux.Router
}

type PlaidClient struct {
	ClientID     string
	Secret       string
	Env          string
	Products     string
	CountryCodes string
	ApiClient    *plaid.APIClient
}

var PlaidEnvs = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}
