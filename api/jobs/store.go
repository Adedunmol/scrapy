package jobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Store interface {
	CreateJob(ctx context.Context, body *CreateJobBody) (Job, error)
	BatchCreateJobs(ctx context.Context, jobs []CreateJobBody) error
	GetJobs(ctx context.Context, userID uuid.UUID, page, limit int) ([]Job, error)
	GetUserCompany(ctx context.Context, userID uuid.UUID) (Company, error)
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

	query := `INSERT INTO jobs (job_title, job_link, date_posted, category_id, origin, origin_id) 
				VALUES (@title, @link, @datePosted, @categoryID, @origin, @originID) 
				RETURNING id, job_title, job_link, date_posted, category_id, origin, origin_id;`
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

	err := row.Scan(&job.ID, &job.JobTitle, &job.JobLink, &job.DatePosted, &job.CategoryID, &job.Origin, &job.OriginID)

	if err != nil {
		err = errors.Join(helpers.ErrInternalServer, err)
		return Job{}, fmt.Errorf("error scanning row (create job): %w", err)
	}

	return job, nil
}

func (j *JobStore) BatchCreateJobs(ctx context.Context, jobs []CreateJobBody) error {
	ctx, cancel := j.WithTimeout(ctx)
	defer cancel()

	// origin_id,
	// @originID,
	query := `
    INSERT INTO jobs (
        job_title,
        job_link,
        date_posted,
        category_id,
        origin,
        updated_at
    ) VALUES (
        @jobTitle,
        @jobLink,
        @datePosted,
        @categoryID,
        @origin,
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
			//"originID":   job.OriginID,
		}
		batch.Queue(query, args)
	}

	results := j.db.SendBatch(ctx, batch)

	defer results.Close()

	for _, _ = range jobs {
		_, err := results.Exec()
		if err != nil {
			log.Println(fmt.Errorf("error while batch creating jobs: %v", err))
			return fmt.Errorf("error while batch creating jobs: %v", err)
		}
	}

	return results.Close()
}

func (j *JobStore) GetJobs(ctx context.Context, userID uuid.UUID, page, limit int) ([]Job, error) {
	ctx, cancel := j.WithTimeout(ctx)
	defer cancel()

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20 // default page size
	}
	offset := (page - 1) * limit

	// add index on preferences.user_id
	// add desc index on created_at
	query := `
			SELECT 
				j.id AS job_id,
				j.job_title,
				j.job_link,
				j.date_posted,
				j.origin,
				j.origin_id,
				j.category_id,
				j.created_at AS job_created_at
			FROM preferences p
			JOIN categories c ON p.categories = c.id
			JOIN jobs j ON j.category_id = c.id
			WHERE p.user_id = @userID
			ORDER BY j.created_at DESC
			LIMIT @limit OFFSET @offset;
			`
	args := pgx.NamedArgs{
		"userID": userID,
		"limit":  limit,
		"offset": offset,
	}

	var jobsData []Job

	rows, err := j.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var job Job

		if err := rows.Scan(&job.ID, &job.JobTitle, &job.JobLink, &job.DatePosted, &job.Origin, &job.OriginID, &job.CategoryID, &job.CreatedAt); err != nil {
			return nil, err
		}

		jobsData = append(jobsData, job)
	}

	return jobsData, nil
}

func (j *JobStore) GetUserCompany(ctx context.Context, userID uuid.UUID) (Company, error) {
	ctx, cancel := j.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, name FROM companies WHERE user_id = @id;"

	args := pgx.NamedArgs{
		"id": userID,
	}

	var company Company

	row := j.db.QueryRow(ctx, query, args)

	err := row.Scan(&company.ID, &company.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Company{}, helpers.ErrNotFound
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Company{}, fmt.Errorf("error scanning row (find user company): %w", err)
	}

	return company, nil
}

func (j *JobStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, j.queryTimeout)
}
