package app

import (
	"net/http"

	"github.com/plaid/plaid-go/plaid"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Env               string
	Port              string
	Secret            string
	JwtKey            string
	RefreshTokenExp   int
	AccessTokenExp    int
	MongoURI          string
	Db                string
	Coll              *Collections
	MongoClient       *mongo.Client
	PlaidClient       *plaid.APIClient
	PlaidClientID     string
	PlaidSecret       string
	PlaidEnv          string
	PlaidCountryCodes string
	PlaidProducts     string
	Server            *http.Server
	Handlers          *Handlers
}

type Handlers struct {
	// Health
	Health http.HandlerFunc
	// Finances
	GetExpenses http.HandlerFunc
	// Plaid
	CreateLinkToken http.HandlerFunc
	SetAccessToken  http.HandlerFunc
	// Users
	CreateUser  http.HandlerFunc
	GetUserInfo http.HandlerFunc
	LoginEmail  http.HandlerFunc
	UserInfo    http.HandlerFunc
	// Web
	ServeClient http.HandlerFunc
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
