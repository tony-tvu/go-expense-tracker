package app

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	// Env               string
	// Port              string
	// EncryptionKey     string
	// JwtKey            string
	// RefreshTokenExp   int
	// AccessTokenExp    int
	// MongoURI          string
	// DbName            string
	Users             *mongo.Collection
	Sessions          *mongo.Collection
	// PlaidClient       *plaid.APIClient
	// PlaidClientID     string
	// PlaidSecret       string
	// PlaidEnv          string
	// PlaidCountryCodes string
	// PlaidProducts     string
	Router            *gin.Engine
}
