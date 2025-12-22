package main

import (
	"log"
	"net/http"

	"github.com/Mx-Tr/pms-go/internal/config"
	"github.com/Mx-Tr/pms-go/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	cfg   *config.Config
	store *store.Storage
}

// Mount настраивает роутер и все хендлеры
func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Префикс /api/v1 - стандарт для REST API
	r.Route("/api/v1", func(r chi.Router) {
		// Публичные руты
		r.Post("/auth/register", app.RegisterUserHandler)
		r.Post("/auth/login", app.LoginUserHandler)

		// Приватные руты
		r.Group(func(r chi.Router) {
			r.Use(app.AuthMiddleware)

			r.Get("/users/me", app.GetCurrentUserHandler)

			r.Get("/projects/getAll", app.GetProjectsHandler)
			r.Post("/projects/create", app.CreateProjectHandler)
			r.Patch("/projects/update/{id}", app.UpdateProjectHandler)
			r.Delete("/projects/delete/{id}", app.DeleteProjectHandler)

			r.Post("/projects/{projectId}/tasks", app.CreateTaskHandler)
		})
	})

	return r
}

// Run запускает сервер
func (app *Application) Run(mux http.Handler) error {
	log.Printf("Server is running on %s", app.cfg.Port)
	return http.ListenAndServe(app.cfg.Port, mux)
}
