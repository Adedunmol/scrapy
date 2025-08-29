package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()

	fmt.Println("registering route")
	r.Post("/fetch-jobs", FetchJobsHandler)
	fmt.Println("registered route")

	return r
}
