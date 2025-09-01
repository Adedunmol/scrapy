package auth

import "github.com/go-chi/chi/v5"

func SetupRoutes(r *chi.Mux) {

	handler := Handler{}

	r.Post("/register", handler.RegisterUserHandler)
	r.Post("/login", handler.LoginUserHandler)
}
