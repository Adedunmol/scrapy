package tests

import (
	"context"
	"errors"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/Adedunmol/scrapy/api/jobs"
	"github.com/google/uuid"
)

type StubJobStore struct {
	Jobs      []jobs.Job
	Fail      bool
	NotFound  bool
	Companies []jobs.Company
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

func (s *StubJobStore) GetUserCompany(ctx context.Context, userID uuid.UUID) (jobs.Company, error) {
	if s.Fail {
		return jobs.Company{}, errors.New("failed to fetch company")
	}
	if s.NotFound {
		return jobs.Company{}, helpers.ErrNotFound
	}

	for _, company := range s.Companies {
		if company.UserID == userID {
			return company, nil
		}
	}
	return jobs.Company{}, helpers.ErrNotFound
}
