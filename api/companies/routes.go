package companies

import (
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func SetupRoutes(r *chi.Mux, db *pgxpool.Pool) {

	companyRouter := chi.NewRouter()
	handler := Handler{
		Store: NewCompanyStore(db, 5*time.Second),
	}

	companyRouter.Use(helpers.AuthMiddleware)
	companyRouter.Post("/", handler.CreateCompany)
	//companyRouter.Get("/", handler.GetCategories)

	r.Mount("/companies", companyRouter)
}
