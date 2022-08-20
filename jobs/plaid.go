package jobs

import (
	"github.com/plaid/plaid-go/plaid"
)

// Plaid Terminology:
// Item = a connection between the user and their bank account

type PlaidClient struct {
	PlaidClientID     string
	PlaidSecret       string
	PlaidEnv          string
	PlaidProducts     string
	PlaidCountryCodes string
	Client            *plaid.APIClient
}

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}
