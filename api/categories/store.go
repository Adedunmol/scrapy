package categories

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store interface {
	CreateCategory(ctx context.Context, body *CreateCategoryBody) (Category, error)
	GetCategories(ctx context.Context) ([]Category, error)
}

type CategoryStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewCategoryStore(db *pgxpool.Pool, queryTimeout time.Duration) *CategoryStore {

	return &CategoryStore{db: db, queryTimeout: queryTimeout}
}

func (c *CategoryStore) CreateCategory(ctx context.Context, body *CreateCategoryBody) (Category, error) {
	ctx, cancel := c.WithTimeout(ctx)
	defer cancel()

	query := `
		INSERT INTO categories (name, created_at, updated_at)
		VALUES (@name, NOW())
		RETURNING id, name, created_at, updated_at;
	`
	args := pgx.NamedArgs{
		"name": body.Name,
	}
	
	var category Category
	err := c.db.QueryRow(ctx, query, args).
		Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func (c *CategoryStore) GetCategories(ctx context.Context) ([]Category, error) {
	ctx, cancel := c.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM categories;"

	var categories []Category

	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return categories, err
	}

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (c *CategoryStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, c.queryTimeout)
}
