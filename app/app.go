package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/plaid/plaid-go/plaid"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Env               string
	Port              string
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
	Handlers          *Handlers
	Router            *mux.Router
}

type Handlers struct {
	// Auth
	Login       http.HandlerFunc
	Logout      http.HandlerFunc
	GetSessions http.HandlerFunc
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

var PlaidEnvs = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}
