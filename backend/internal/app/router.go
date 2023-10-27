package app

import (
	mlogger "emtest/internal/middleware/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *App) NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*",
			"http://*",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestID)
	r.Use(mlogger.New(app.Log))

	r.Get("/persons", app.FilterPersons)
	r.Post("/person", app.StorePerson)
	r.Delete("/person/{id}", app.DeletePersonById)
	r.Put("/person/{id}", app.UpdatePersonById)
	return r
}
