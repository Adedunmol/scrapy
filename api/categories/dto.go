package categories

import (
	"github.com/google/uuid"
	"time"
)

type CreateCategoryBody struct {
	Name string `json:"name" validate:"required"`
}

type Category struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
