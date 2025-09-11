package companies

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/Adedunmol/scrapy/api/jobs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store interface {
	GetUserCompany(ctx context.Context, userID uuid.UUID) (Company, error)
	CreateCompany(ctx context.Context, body *CreateCompanyBody) (Company, error)
	GetCompany(ctx context.Context, companyID uuid.UUID) (Company, error)
	GetCompanyJobs(ctx context.Context, companyID uuid.UUID) ([]jobs.Job, error)
}

const UniqueViolationCode = "23505"

type CompanyStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewCompanyStore(db *pgxpool.Pool, queryTimeout time.Duration) *CompanyStore {

	return &CompanyStore{db: db, queryTimeout: queryTimeout}
}

func (c *CompanyStore) GetUserCompany(ctx context.Context, userID uuid.UUID) (Company, error) {
	ctx, cancel := c.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, name, email, created_at updated_at FROM companies WHERE user_id = @id;"

	args := pgx.NamedArgs{
		"id": userID,
	}

	var company Company

	row := c.db.QueryRow(ctx, query, args)

	err := row.Scan(&company.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Company{}, helpers.ErrNotFound
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Company{}, fmt.Errorf("error scanning row (find user company): %w", err)
	}

	return company, nil
}

func (c *CompanyStore) CreateCompany(ctx context.Context, body *CreateCompanyBody) (Company, error) {
	ctx, cancel := c.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO companies (name, email, user_id) VALUES (@name, @email, @userID) RETURNING id, name, email, user_id, created_at, updated_at;"
	args := pgx.NamedArgs{
		"email":  body.Email,
		"name":   body.Name,
		"userID": body.UserID,
	}

	var company Company

	row := c.db.QueryRow(ctx, query, args)

	err := row.Scan(&company.ID, &company.Name, &company.Email, &company.UserID, &company.CreatedAt, &company.UpdatedAt)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == UniqueViolationCode {
			return Company{}, helpers.ErrConflict
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Company{}, fmt.Errorf("error scanning row (create company): %w", err)
	}

	return company, nil
}

func (c *CompanyStore) GetCompany(ctx context.Context, companyID uuid.UUID) (Company, error) {
	ctx, cancel := c.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, name, email, created_at updated_at FROM companies WHERE id = @id;"

	args := pgx.NamedArgs{
		"id": companyID,
	}

	var company Company

	row := c.db.QueryRow(ctx, query, args)

	err := row.Scan(&company.ID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Company{}, helpers.ErrNotFound
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Company{}, fmt.Errorf("error scanning row (find company): %w", err)
	}

	return company, nil
}

func (c *CompanyStore) GetCompanyJobs(ctx context.Context, companyID uuid.UUID) ([]jobs.Job, error) {
	ctx, cancel := c.WithTimeout(ctx)
	defer cancel()

	// add index on origin_id
	query := `SELECT id, job_title, job_link, date_posted, origin, origin_id, category_id, created_at FROM jobs WHERE origin_id = @companyID;`

	args := pgx.NamedArgs{
		"companyID": companyID,
	}

	var jobsData []jobs.Job

	rows, err := c.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var job jobs.Job

		if err := rows.Scan(&job.ID, &job.JobTitle, &job.JobLink, &job.DatePosted, &job.Origin, &job.OriginID, &job.CategoryID, &job.CreatedAt); err != nil {
			return nil, err
		}

		jobsData = append(jobsData, job)
	}

	return jobsData, nil
}

func (c *CompanyStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, c.queryTimeout)
}
