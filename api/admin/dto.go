package admin

import (
	"github.com/google/uuid"
	"time"
)

type CreateAdminBody struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required"`
	//Username    string   `json:"username" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type Admin struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
}

type CreateRoleBody struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Slug        string `json:"slug"`
	CreatedBy   string `json:"created_by"`
}

type Role struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Slug        string     `json:"slug"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type CreatePermissionBody struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Slug        string `json:"slug"`
	CreatedBy   string `json:"created_by"`
}

type Permission struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Slug        string     `json:"slug"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
