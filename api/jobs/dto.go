package jobs

import (
	"github.com/google/uuid"
	"time"
)

type CreateJobBody struct {
	JobTitle    string    `json:"job_title" validate:"required"`
	JobLink     string    `json:"job_link" validate:"required"`
	DatePosted  string    `json:"date_posted" validate:"required"`
	Category    string    `json:"category" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Origin      string    `json:"origin"`
	OriginID    uuid.UUID `json:"origin_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	CompanyName string    `json:"company_name"`
}

type Job struct {
	ID           uuid.UUID  `json:"id"`
	JobTitle     string     `json:"job_title"`
	JobLink      string     `json:"job_link"`
	Description  string     `json:"description"`
	DatePosted   string     `json:"date_posted"`
	CategoryID   uuid.UUID  `json:"category_id"`
	Origin       string     `json:"origin"`
	OriginID     uuid.UUID  `json:"origin_id,omitempty"`
	CategoryName string     `json:"category_name"`
	CompanyName  string     `json:"company_name"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type Company struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id"`
}
