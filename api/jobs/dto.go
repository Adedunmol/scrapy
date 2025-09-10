package jobs

import "github.com/google/uuid"

type CreateJobBody struct {
	JobTitle   string    `json:"job_title" validate:"required"`
	JobLink    string    `json:"job_link" validate:"required"`
	DatePosted string    `json:"date_posted" validate:"required"`
	Category   string    `json:"category" validate:"required"`
	Origin     string    `json:"origin"`
	OriginID   uuid.UUID `json:"origin_id"`
	CategoryID uuid.UUID `json:"category_id"`
}

type Job struct {
	ID         uuid.UUID `json:"id"`
	JobTitle   string    `json:"job_title"`
	JobLink    string    `json:"job_link"`
	DatePosted string    `json:"date_posted"`
	CategoryID uuid.UUID `json:"category_id"`
	Origin     string    `json:"origin"`
	OriginID   uuid.UUID `json:"origin_id"`
}
