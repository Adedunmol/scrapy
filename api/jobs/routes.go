package jobs

import (
	"github.com/Adedunmol/scrapy/api/categories"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/Adedunmol/scrapy/api/transactions"
	"github.com/Adedunmol/scrapy/api/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func SetupRoutes(r *chi.Mux, db *pgxpool.Pool) {

	jobRouter := chi.NewRouter()

	handler := Handler{
		Store:            NewJobStore(db, 5*time.Second),
		CategoriesStore:  categories.NewCategoryStore(db, 5*time.Second),
		WalletStore:      wallet.NewWalletStore(db, 5*time.Second),
		TransactionStore: transactions.NewTransactionStore(db, 5*time.Second),
	}

	jobRouter.Use(helpers.AuthMiddleware)

	jobRouter.Post("/", handler.CreateJobHandler)
	jobRouter.Get("/", handler.GetUserJobsHandler)

	r.Mount("/jobs", jobRouter)
}
