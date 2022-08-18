package handlers

import (
	"fmt"
	"net/http"

	"github.com/tony-tvu/goexpense/app"
)

func GetExpenses(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "GetExpenses called")
	}
}
