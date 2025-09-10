package tests

import (
	"context"
	"errors"
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/categories"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
)

type StubCategoryStore struct {
	Users      []auth.User
	Fail       bool
	Conflict   bool
	Companies  []auth.Company
	Categories []categories.Category
}

const TestCategoryID = "123e4567-e89b-12d3-a456-426614174000"

func (c *StubCategoryStore) CreateCategory(ctx context.Context, body *categories.CreateCategoryBody) (categories.Category, error) {
	if c.Fail {
		return categories.Category{}, errors.New("store error")
	}
	if c.Conflict {
		return categories.Category{}, helpers.ErrConflict
	}

	category := categories.Category{
		ID:   uuid.MustParse(TestCategoryID),
		Name: body.Name,
	}

	c.Categories = append(c.Categories, category)
	return category, nil
}

func (c *StubCategoryStore) GetCategories(ctx context.Context) ([]categories.Category, error) {
	if c.Fail {
		return nil, errors.New("store error")
	}
	return c.Categories, nil
}
func (c *StubCategoryStore) GetCategory(ctx context.Context, name string) (categories.Category, error) {
	if c.Fail {
		return categories.Category{}, errors.New("store error")
	}

	for _, cat := range c.Categories {
		if cat.Name == name {
			return cat, nil
		}
	}
	return categories.Category{}, helpers.ErrNotFound
}
