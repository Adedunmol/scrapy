package api

import (
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func Routes(db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // frontends allowed
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by browsers
	}))

	r.Get("/check", func(w http.ResponseWriter, r *http.Request) {

		helpers.WriteJSONResponse(w, "hello from scrapy", http.StatusOK)
	})

	r.Post("/fetch-jobs", FetchJobsHandler)

	auth.SetupRoutes(r, db)

	return r
}
