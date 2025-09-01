package tests

import (
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/helpers"
)

type StubUserStore struct {
	Users []auth.User
	Fail  bool
}

func (s *StubUserStore) CreateUser(body *auth.CreateUserBody) error {

	if !s.Fail {
		for _, u := range s.Users {
			if u.Email == body.Email {
				return helpers.ErrConflict
			}
		}

		userData := auth.User{ID: 1, FirstName: body.FirstName, LastName: body.LastName, Username: body.Username, Email: body.Email, Password: body.Password}

		s.Users = append(s.Users, userData)

		return nil
	}

	return helpers.ErrInternalServer
}

func (s *StubUserStore) FindUserByEmail(email string) (auth.User, error) {

	for _, u := range s.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return auth.User{}, helpers.ErrNotFound
}

func (s *StubUserStore) ComparePasswords(storedPassword, candidatePassword string) bool {
	return storedPassword == candidatePassword
}
