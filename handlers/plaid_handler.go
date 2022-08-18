package handlers

import (
	"fmt"
	"net/http"

	"github.com/tony-tvu/goexpense/app"
)

func PlaidHandler(a *app.App) func(w http.ResponseWriter, r *http.Request) {
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
