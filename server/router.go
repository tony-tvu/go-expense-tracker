package server

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/handlers"
)

func InitRouter(ctx context.Context, a *app.App) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/health",
		Chain(handlers.Health, Middlewares...)).Methods("GET")
	router.Handle("/api/login_email",
		Chain(handlers.LoginEmail(ctx, a), Logging(), LoginRateLimit(), NoCache())).Methods("POST")
	router.Handle("/api/get_token_exp",
		Chain(handlers.IsTokenValid(ctx, a), Logging(), LoginRateLimit(), NoCache())).Methods("POST")
	router.Handle("/api/user",
		Chain(handlers.CreateUser(ctx, a), Middlewares...)).Methods("POST")
	router.Handle("/api/expense",
		Chain(handlers.GetExpenses(ctx, a), Middlewares...)).Methods("GET")
	// TODO: add auth to this so only registered users can create link tokens
	router.Handle("/api/create_link_token",
		Chain(handlers.CreateLinkToken(ctx, a), Middlewares...)).Methods("GET")
	router.Handle("/api/set_access_token",
		Chain(handlers.SetAccessToken(ctx, a), Middlewares...)).Methods("POST")
	router.PathPrefix("/").Handler(
		Chain(handlers.SpaHandler("web/build", "index.html"), Middlewares...)).Methods("GET")
	return router
}
