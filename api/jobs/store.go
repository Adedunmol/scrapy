package jobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store interface {
	CreateJob(ctx context.Context, body *CreateJobBody) (Job, error)
	BatchCreateJobs(ctx context.Context, jobs []CreateJobBody) error
}

type JobStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewJobStore(db *pgxpool.Pool, queryTimeout time.Duration) *JobStore {

	return &JobStore{db: db, queryTimeout: queryTimeout}
}

func (j *JobStore) CreateJob(ctx context.Context, body *CreateJobBody) (Job, error) {
	ctx, cancel := j.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO jobs (job_title, job_link, date_posted, category_id, origin, origin_id) VALUES (@title, @link, @datePosted, @categoryID, @origin, @originID) RETURNING id, job_title, job_link, date_posted, category_id, origin, origin_id;"
	args := pgx.NamedArgs{
		"title":      body.JobTitle,
		"link":       body.JobLink,
		"datePosted": body.DatePosted,
		"categoryID": body.CategoryID,
		"origin":     body.Origin,
		"originID":   body.OriginID,
	}

	var job Job

	row := j.db.QueryRow(ctx, query, args)

	err := row.Scan(&job.ID, &job.JobTitle, &job.JobLink, &job.DatePosted, &job.CategoryID, &job.Origin)

	if err != nil {
		err = errors.Join(helpers.ErrInternalServer, err)
		return Job{}, fmt.Errorf("error scanning row (create user): %w", err)
	}

	return job, nil
}

func (j *JobStore) BatchCreateJobs(ctx context.Context, jobs []CreateJobBody) error {
	ctx, cancel := j.WithTimeout(ctx)
	defer cancel()

	query := `
    INSERT INTO jobs (
        job_title,
        job_link,
        date_posted,
        category_id,
        origin,
        origin_id,
        updated_at
    ) VALUES (
        @jobTitle,
        @jobLink,
        @datePosted,
        @categoryID,
        @origin,
        @originID,
        NOW()
    )`

	batch := &pgx.Batch{}

	for _, job := range jobs {
		args := pgx.NamedArgs{
			"jobTitle":   job.JobTitle,
			"jobLink":    job.JobLink,
			"datePosted": job.DatePosted,
			"categoryID": job.CategoryID,
			"origin":     job.Origin,
			"originID":   job.OriginID,
		}
		batch.Queue(query, args)
	}

	results := j.db.SendBatch(ctx, batch)
	defer results.Close()

	for _, _ = range jobs {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("error while creating preferences: %w", err)
		}
	}

	return results.Close()
}

func (j *JobStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, j.queryTimeout)
}
