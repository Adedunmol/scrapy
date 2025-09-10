package companies

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store interface {
	GetUserCompany(ctx context.Context, userID uuid.UUID) (Company, error)
}

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

func (c *CompanyStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, c.queryTimeout)
}
