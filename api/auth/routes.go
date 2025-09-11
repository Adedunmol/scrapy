package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func SetupRoutes(r *chi.Mux, db *pgxpool.Pool) {

	handler := Handler{
		Store: NewUserStore(db, 5*time.Second),
	}

	r.Post("/register", handler.RegisterUserHandler)
	r.Post("/login", handler.LoginUserHandler)
}
