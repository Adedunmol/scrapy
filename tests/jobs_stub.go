package tests

import (
	"context"
	"errors"
	"github.com/Adedunmol/scrapy/api/jobs"
	"github.com/google/uuid"
)

type StubJobStore struct {
	Jobs []jobs.Job
	Fail bool
}

func (s *StubJobStore) CreateJob(ctx context.Context, body *jobs.CreateJobBody) (jobs.Job, error) {
	if s.Fail {
		return jobs.Job{}, errors.New("failed to create job")
	}

	job := jobs.Job{
		ID:         uuid.New(),
		JobTitle:   body.JobTitle,
		JobLink:    body.JobLink,
		CategoryID: body.CategoryID,
		Origin:     body.Origin,
		OriginID:   body.OriginID,
	}
	s.Jobs = append(s.Jobs, job)
	return job, nil
}

func (s *StubJobStore) BatchCreateJobs(ctx context.Context, jobs []jobs.CreateJobBody) error {
	return nil
}

func (s *StubJobStore) GetJobs(ctx context.Context, userID uuid.UUID) ([]jobs.Job, error) {
	return nil, nil
}
