package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store interface {
	CreateUser(body *CreateUserBody) (User, error)
	FindUserByEmail(email string) (User, error)
	ComparePasswords(password, candidatePassword string) bool
	GetCategories(ctx context.Context) (map[string]uuid.UUID, error)
	CreatePreferences(ctx context.Context, preferences []uuid.UUID, userID uuid.UUID) error
}

type UserStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewAuthStore(db *pgxpool.Pool, queryTimeout time.Duration) *UserStore {

	return &UserStore{db: db, queryTimeout: queryTimeout}
}

func (s *UserStore) CreateUser(body *CreateUserBody) (User, error) {

	return User{}, nil
}

func (s *UserStore) FindUserByEmail(email string) (User, error) {

	return User{}, nil
}

func (s *UserStore) ComparePasswords(password, candidatePassword string) bool {
	return false
}

func (s *UserStore) CreatePreferences(ctx context.Context, preferences []int, userID uuid.UUID) error {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO preferences (user_id, catgories) VALUES (@userID, @categoryID)"
	batch := &pgx.Batch{}

	for _, pref := range preferences {
		args := pgx.NamedArgs{
			"userID":     userID,
			"categoryID": pref,
		}
		batch.Queue(query, args)
	}

	results := s.db.SendBatch(ctx, batch)
	defer results.Close()

	for _, _ = range preferences {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("error while creating preferences: %w", err)
		}
	}

	return results.Close()
}

func (s *UserStore) GetCategories(ctx context.Context) (map[string]uuid.UUID, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, name FROM categories"

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	categories := make(map[string]uuid.UUID)
	for rows.Next() {
		var id uuid.UUID
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		categories[name] = id
	}

	return categories, nil
}

func (s *UserStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, s.queryTimeout)
}
