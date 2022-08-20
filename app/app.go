package app

import (
	"github.com/gorilla/mux"
	"github.com/plaid/plaid-go/plaid"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Env             string
	Port            string
	Secret          string
	JwtKey          string
	RefreshTokenExp int
	AccessTokenExp  int
	MongoURI        string
	Db              string
	Coll            *Collections
	MongoClient     *mongo.Client
	PlaidClient     *plaid.APIClient
	CountryCodes    string
	Products        string
	Router          *mux.Router
}

type Collections struct {
	Users    string
	Sessions string
}

var PlaidEnvs = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}
