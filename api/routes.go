package api

import (
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Routes(db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/fetch-jobs", FetchJobsHandler)

	auth.SetupRoutes(r, db)

	return r
}
