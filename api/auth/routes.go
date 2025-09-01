package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *chi.Mux, db *pgxpool.Pool) {

	handler := Handler{}

	r.Post("/register", handler.RegisterUserHandler)
	r.Post("/login", handler.LoginUserHandler)
}
