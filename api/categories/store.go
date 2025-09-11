package categories

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
	CreateCategory(ctx context.Context, body *CreateCategoryBody) (Category, error)
	GetCategories(ctx context.Context) ([]Category, error)
	GetCategory(ctx context.Context, name string) (Category, error)
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
		INSERT INTO categories (name, updated_at)
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
		return Category{}, fmt.Errorf("error inserting category: %v", err)
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
		return categories, fmt.Errorf("error fetching categories: %v", err)
	}

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning categories: %v", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (c *CategoryStore) GetCategory(ctx context.Context, name string) (Category, error) {
	ctx, cancel := c.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, name, created_at, updated_at FROM categories WHERE name = @name;"

	args := pgx.NamedArgs{
		"name": name,
	}

	var category Category

	row := c.db.QueryRow(ctx, query, args)

	err := row.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Category{}, helpers.ErrNotFound
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return Category{}, fmt.Errorf("error scanning row (find category): %w", err)
	}

	return category, nil
}

func (c *CategoryStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, c.queryTimeout)
}
