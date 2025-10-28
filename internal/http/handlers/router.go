package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/medidew/ApplicationTracker/internal/http/middleware"
)

func SetupRouter(app *App) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.ZapLoggerMiddleware(app.Logger))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

	})

	router.Route("/applications", func(router chi.Router) {
		router.Get("/", app.ListApplications)
		router.Post("/", app.CreateApplication)

		router.Route("/{companyID}", func(router chi.Router) {
			router.Get("/", app.GetApplication)
			router.Delete("/", app.DeleteApplication)
			router.Put("/", app.UpdateApplicationStatus)

			router.Route("/notes", func(router chi.Router) {
				router.Get("/", app.ListApplicationNotes)
				router.Post("/", app.AddApplicationNote)
				router.Delete("/{noteIndex}", app.RemoveApplicationNote)
			})
		})
	})

	router.Get("/register", app.Register)
	router.Get("/login", app.Login)
	router.Get("/logout", app.Logout)

	return router
}