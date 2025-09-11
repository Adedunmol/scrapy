package tests

import (
	"context"
	"errors"
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/companies"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
)

type StubCompanyStore struct {
	Companies []companies.Company
	Users     []auth.User

	Fail     bool // simulate generic failure
	NotFound bool // simulate "not found"
	Conflict bool // simulate duplicate company conflict
}

func (s *StubCompanyStore) GetUserCompany(ctx context.Context, userID uuid.UUID) (companies.Company, error) {
	if s.Fail {
		return companies.Company{}, errors.New("failed to fetch company")
	}
	if s.NotFound {
		return companies.Company{}, helpers.ErrNotFound
	}

	for _, company := range s.Companies {
		if company.UserID == userID {
			return company, nil
		}
	}
	return companies.Company{}, helpers.ErrNotFound
}

func (s *StubCompanyStore) CreateCompany(ctx context.Context, body *companies.CreateCompanyBody) (companies.Company, error) {
	if s.Fail {
		return companies.Company{}, helpers.ErrInternalServer
	}

	if s.Conflict {
		return companies.Company{}, helpers.ErrConflict
	}

	for _, c := range s.Companies {
		if c.Email == body.Email {
			return companies.Company{}, helpers.ErrConflict
		}
	}

	newCompany := companies.Company{
		ID:     uuid.New(),
		Email:  body.Email,
		Name:   body.Name,
		UserID: body.UserID,
	}

	s.Companies = append(s.Companies, newCompany)
	return newCompany, nil
}
