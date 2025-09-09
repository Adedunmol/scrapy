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
	"golang.org/x/crypto/bcrypt"
	"time"
)

const UniqueViolationCode = "23505"

type Store interface {
	CreateUser(ctx context.Context, body *CreateUserBody) (User, error)
	FindUserByEmail(ctx context.Context, email string) (User, error)
	ComparePasswords(password, candidatePassword string) bool
	GetCategories(ctx context.Context) (map[string]uuid.UUID, error)
	CreatePreferences(ctx context.Context, preferences []uuid.UUID, userID uuid.UUID) error
	CreateCompany(ctx context.Context, body *CreateCompanyBody) (Company, error)
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

	query := "INSERT INTO users (email, first_name, last_name, password) VALUES (@email, @firstName, @lastName, @password) RETURNING id, email, first_name, last_name;"
	args := pgx.NamedArgs{
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

func (s *UserStore) FindUserByEmail(ctx context.Context, email string) (User, error) {

	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "SELECT id, first_name, last_name, email, password FROM users WHERE email = @email;"
	args := pgx.NamedArgs{
		"email": email,
	}

	var user User

	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password)

	if err != nil {
		err = errors.Join(helpers.ErrInternalServer, err)
		return User{}, fmt.Errorf("error scanning row (find user by email): %w", err)
	}

	return user, nil
}

func (s *UserStore) ComparePasswords(password, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(candidatePassword))

	if err != nil {
		return false
	}
	return true
}

func (s *UserStore) CreatePreferences(ctx context.Context, preferences []uuid.UUID, userID uuid.UUID) error {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO preferences (user_id, categories) VALUES (@userID, @categoryID)"
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

func (s *UserStore) CreateCompany(ctx context.Context, body *CreateCompanyBody) (Company, error) {
	ctx, cancel := s.WithTimeout(ctx)
	defer cancel()

	query := "INSERT INTO companies (name, email, user_id) VALUES (@name, @email, @userID) RETURNING id, name, email, user_id, created_at, updated_at;"
	args := pgx.NamedArgs{
		"email":  body.Email,
		"name":   body.Name,
		"userID": body.UserID,
	}

	var company Company

	row := s.db.QueryRow(ctx, query, args)

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

func (s *UserStore) WithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, s.queryTimeout)
}
