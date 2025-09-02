package api

import (
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func Routes(db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/check", func(w http.ResponseWriter, r *http.Request) {

		helpers.WriteJSONResponse(w, "hello from scrapy", http.StatusOK)
	})

	r.Post("/fetch-jobs", FetchJobsHandler)

	auth.SetupRoutes(r, db)

	return r
}
