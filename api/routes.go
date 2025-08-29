package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/fetch-jobs", http.HandlerFunc(FetchJobsHandler))

	return r
}
