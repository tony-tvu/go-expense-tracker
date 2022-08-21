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
		Chain(handlers.Health, LoginProtected(ctx, a))).Methods("GET")
	router.Handle("/api/login_email",
		Chain(handlers.LoginEmail(ctx, a), UseMiddlewares(LoginRateLimit())...)).Methods("POST")
	router.Handle("/api/user",
		Chain(handlers.CreateUser(ctx, a), UseMiddlewares()...)).Methods("POST")
	router.Handle("/api/expense",
		Chain(handlers.GetExpenses(ctx, a), UseMiddlewares()...)).Methods("GET")
	// TODO: add auth to this so only registered users can create link tokens
	router.Handle("/api/create_link_token",
		Chain(handlers.CreateLinkToken(ctx, a), UseMiddlewares()...)).Methods("GET")
	router.Handle("/api/set_access_token",
		Chain(handlers.SetAccessToken(ctx, a), UseMiddlewares()...)).Methods("POST")
	router.Handle("/api/admin",
		Chain(handlers.GetTokens(ctx, a), UseMiddlewares()...)).Methods("GET")
	router.PathPrefix("/").Handler(
		Chain(handlers.SpaHandler("web/build", "index.html"), UseMiddlewares()...)).Methods("GET")
	return router
}
