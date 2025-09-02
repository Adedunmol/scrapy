package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const UniqueViolationCode = "23505"

type Store interface {
	CreateUser(ctx context.Context, body *CreateUserBody) (User, error)
	FindUserByEmail(email string) (User, error)
	ComparePasswords(password, candidatePassword string) bool
	GetCategories(ctx context.Context) (map[string]uuid.UUID, error)
	CreatePreferences(ctx context.Context, preferences []uuid.UUID, userID uuid.UUID) error
}

type UserStore struct {
	db           *pgxpool.Pool
	queryTimeout time.Duration
}

func NewUserStore(db *pgxpool.Pool, queryTimeout time.Duration) *UserStore {

	return &UserStore{db: db, queryTimeout: queryTimeout}
}

func (s *UserStore) CreateUser(ctx context.Context, body *CreateUserBody) (User, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO users (name, email, first_name, last_name, password) VALUES (@username, @email, @firstName, @lastName, @password) RETURNING id, email, first_name, last_name;"
	args := pgx.NamedArgs{
		"username":  body.Username,
		"email":     body.Email,
		"firstName": body.FirstName,
		"lastName":  body.LastName,
		"password":  body.Password,
	}

	var user User

	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName)

	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == UniqueViolationCode {
			return User{}, helpers.ErrConflict
		}
		err = errors.Join(helpers.ErrInternalServer, err)
		return User{}, fmt.Errorf("error scanning row (create user): %w", err)
	}

	return user, nil
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
