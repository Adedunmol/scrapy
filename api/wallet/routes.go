package wallet

import (
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func SetupRoutes(r *chi.Mux, db *pgxpool.Pool) {

	jobRouter := chi.NewRouter()

	handler := Handler{
		Store: NewWalletStore(db, 5*time.Second),
	}

	jobRouter.Use(helpers.AuthMiddleware)

	jobRouter.Get("/", handler.GetWalletHandler)
	jobRouter.Patch("/", handler.TopUpWalletHandler)

	r.Mount("/wallets", jobRouter)
}
