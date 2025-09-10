package companies

import (
	"github.com/google/uuid"
	"time"
)

type Company struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	UserID    uuid.UUID  `json:"user_id"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
