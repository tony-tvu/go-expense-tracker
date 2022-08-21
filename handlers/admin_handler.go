package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tony-tvu/goexpense/app"
)

// Admin only protected route
func GetTokens(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprint(w, "Admin-only GetTokens accessed")
	}
}
