package handlers

import (
	"fmt"
	"github.com/tony-tvu/goexpense/config"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func ExpenseHandler(cfg *config.AppConfig, client *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {

			fmt.Fprint(w, "Scrape succeeded")
			return
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}
