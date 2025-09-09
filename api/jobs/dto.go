package jobs

import "github.com/google/uuid"

type CreateJobBody struct {
	JobTitle   string    `json:"job_title"`
	JobLink    string    `json:"job_link"`
	DatePosted string    `json:"date_posted"`
	CategoryID uuid.UUID `json:"category_id"`
	Origin     string    `json:"origin"`
	OriginID   int       `json:"origin_id"`
}

type Job struct {
	ID         uuid.UUID `json:"id"`
	JobTitle   string    `json:"job_title"`
	JobLink    string    `json:"job_link"`
	DatePosted string    `json:"date_posted"`
	CategoryID uuid.UUID `json:"category_id"`
	Origin     string    `json:"origin"`
	OriginID   string    `json:"origin_id"`
}
