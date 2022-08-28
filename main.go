package main

import (
	"context"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/server"
)

func main() {

	// s := &server.Server{
	// 	App: &app.App{
	// 		Env:               os.Getenv("ENV"),
	// 		Port:              os.Getenv("PORT"),
	// 		EncryptionKey:     os.Getenv("ENCRYPTION_KEY"),
	// 		JwtKey:            os.Getenv("JWT_KEY"),
	// 		RefreshTokenExp:   86400,
	// 		AccessTokenExp:    900,
	// 		MongoURI:          os.Getenv("MONGODB_URI"),
	// 		DbName:            os.Getenv("DB_NAME"),
	// 		PlaidClientID:     os.Getenv("PLAID_CLIENT_ID"),
	// 		PlaidSecret:       os.Getenv("PLAID_SECRET"),
	// 		PlaidEnv:          os.Getenv("PLAID_ENV"),
	// 		PlaidCountryCodes: "US,CA",
	// 		PlaidProducts:     "auth,transactions",
	// 	},
	// }

	s := &server.Server{App: &app.App{}}
	s.Initialize()
	s.Run(context.Background())
}
