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
	PlaidClient    *plaid.APIClient
	Router         *mux.Router
}

var PlaidEnvs = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}
