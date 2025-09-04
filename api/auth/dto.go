package auth

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserBody struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required"`
	//Username    string   `json:"username" validate:"required"`
	Email       string   `json:"email" validate:"required,email"`
	SearchTerms []string `json:"search_terms" validate:"required"`
}

type LoginUserBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateCompanyBody struct {
	Name   string    `json:"name" validate:"required"`
	Email  string    `json:"email" validate:"required,email"`
	UserID uuid.UUID `json:"user_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}

type Company struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	UserID    uuid.UUID  `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
