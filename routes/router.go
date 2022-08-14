package routes

import (
	"github.com/gorilla/mux"
	"github.com/tony-tvu/goexpense/config"
	"github.com/tony-tvu/goexpense/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetRouter(cfg *config.AppConfig, client *mongo.Client) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", Chain(handlers.HealthHandler, Middlewares...))
	router.Handle("/api/user", Chain(handlers.UserHandler(cfg, client), Middlewares...))
	router.Handle("/api/expense", Chain(handlers.ExpenseHandler(cfg, client), Middlewares...))
	router.PathPrefix("/").Handler(Chain(handlers.SpaHandler("web/build", "index.html"), Middlewares...))

	return router
}
