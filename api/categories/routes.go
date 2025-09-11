package categories

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func SetupRoutes(r *chi.Mux, db *pgxpool.Pool) {

	categoryRouter := chi.NewRouter()
	handler := Handler{
		Store: NewCategoryStore(db, 5*time.Second),
	}

	categoryRouter.Post("/", handler.CreateCategory) // make it protected
	categoryRouter.Get("/", handler.GetCategories)

	r.Mount("/categories", categoryRouter)
}
