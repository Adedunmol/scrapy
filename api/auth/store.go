package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	CreateUser(body *CreateUserBody) error
	FindUserByEmail(email string) (User, error)
	ComparePasswords(password, candidatePassword string) bool
}

type UserStore struct{ db *pgxpool.Pool }

func NewAuthStore(db *pgxpool.Pool) *UserStore {

	return &UserStore{db: db}
}

func (s *UserStore) CreateUser(body *CreateUserBody) error {

	return nil
}

func (s *UserStore) FindUserByEmail(email string) (User, error) {

	return User{}, nil
}

func (s *UserStore) ComparePasswords(password, candidatePassword string) bool {
	return false
}
