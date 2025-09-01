package api

import (
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/fetch-jobs", FetchJobsHandler)

	auth.SetupRoutes(r)

	return r
}
