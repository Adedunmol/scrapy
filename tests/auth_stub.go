package tests

import (
	"context"
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
)

type StubUserStore struct {
	Users     []auth.User
	Fail      bool
	Companies []auth.Company
}

func (s *StubUserStore) CreateUser(ctx context.Context, body *auth.CreateUserBody) (auth.User, error) {

	if !s.Fail {
		for _, u := range s.Users {
			if u.Email == body.Email {
				return auth.User{}, helpers.ErrConflict
			}
		}

		// ID: 1,
		userData := auth.User{FirstName: body.FirstName, LastName: body.LastName, Email: body.Email, Password: body.Password}

		s.Users = append(s.Users, userData)

		return userData, nil
	}

	return auth.User{}, helpers.ErrInternalServer
}

func (s *StubUserStore) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {

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

func (s *StubUserStore) GetCategories(ctx context.Context) (map[string]uuid.UUID, error) {
	return nil, nil
}

func (s *StubUserStore) CreatePreferences(ctx context.Context, preferences []uuid.UUID, userID uuid.UUID) error {
	return nil
}

func (s *StubUserStore) CreateCompany(ctx context.Context, body *auth.CreateCompanyBody) (auth.Company, error) {

	if !s.Fail {
		for _, u := range s.Companies {
			if u.Email == body.Email {
				return auth.Company{}, helpers.ErrConflict
			}
		}

		// ID: 1,
		companyData := auth.Company{Email: body.Email, Name: body.Name, UserID: body.UserID}

		s.Companies = append(s.Companies, companyData)

		return companyData, nil
	}

	return auth.Company{}, helpers.ErrInternalServer
}
