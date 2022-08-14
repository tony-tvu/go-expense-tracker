package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tony-tvu/goexpense/config"
	"github.com/tony-tvu/goexpense/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Config      *config.AppConfig
	MongoClient *mongo.Client
}

func (app *App) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", Chain(handlers.HealthHandler, Middlewares...))
	router.Handle("/api/user", Chain(handlers.UserHandler(app.Config, app.MongoClient), Middlewares...))
	router.PathPrefix("/").Handler(Chain(handlers.SpaHandler("web/build", "index.html"), Middlewares...))

	log.Printf("Listening on port %s", app.Config.Port)
	if err := http.ListenAndServe(":"+app.Config.Port, router); err != nil {
		log.Fatal(err)
	}
}
