package categories_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Adedunmol/scrapy/api/categories"
	"github.com/Adedunmol/scrapy/tests"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateCategory(t *testing.T) {
	t.Run("successfully creates a category", func(t *testing.T) {
		store := &tests.StubCategoryStore{}
		handler := &categories.Handler{Store: store}

		data := []byte(`{"name": "Engineering"}`)
		req := createCategoryRequest(data)
		rec := httptest.NewRecorder()

		handler.CreateCategory(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		want := map[string]interface{}{
			"status":  "success",
			"message": "category created successfully",
		}

		assertResponseCode(t, rec.Code, http.StatusCreated)

		if len(store.Categories) != 1 {
			t.Errorf("expected 1 category in store, got %d", len(store.Categories))
		}

		if got["message"] != want["message"] {
			t.Errorf("expected message to be %s, got %s", want["message"], got["message"])
		}

		if store.Categories[0].Name != "Engineering" {
			t.Errorf("expected category name to be Engineering, got %s", store.Categories[0].Name)
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		store := &tests.StubCategoryStore{}
		handler := &categories.Handler{Store: store}

		data := []byte(`{"name": "Engineering"`) // missing closing brace
		req := createCategoryRequest(data)
		rec := httptest.NewRecorder()

		handler.CreateCategory(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)
		assertResponseBody(t, got, map[string]interface{}{
			"status":  "error",
			"message": "error decoding body",
		})
	})

	t.Run("validation error returns 400", func(t *testing.T) {
		store := &tests.StubCategoryStore{}
		handler := &categories.Handler{Store: store}

		// Missing required name field
		data := []byte(`{}`)
		req := createCategoryRequest(data)
		rec := httptest.NewRecorder()

		handler.CreateCategory(rec, req)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})

	t.Run("store returns conflict (409)", func(t *testing.T) {
		store := &tests.StubCategoryStore{Conflict: true}
		handler := &categories.Handler{Store: store}

		data := []byte(`{"name": "Engineering"}`)
		req := createCategoryRequest(data)
		rec := httptest.NewRecorder()

		handler.CreateCategory(rec, req)

		assertResponseCode(t, rec.Code, http.StatusConflict)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})

	t.Run("store returns generic error (500)", func(t *testing.T) {
		store := &tests.StubCategoryStore{Fail: true}
		handler := &categories.Handler{Store: store}

		data := []byte(`{"name": "Engineering"}`)
		req := createCategoryRequest(data)
		rec := httptest.NewRecorder()

		handler.CreateCategory(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestGetCategoriesHandler(t *testing.T) {
	t.Run("successfully fetches categories", func(t *testing.T) {
		store := &tests.StubCategoryStore{
			Categories: []categories.Category{
				{Name: "Engineering"},
				{Name: "Marketing"},
			},
		}
		handler := &categories.Handler{Store: store}

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()

		handler.GetCategories(rec, req)

		assertResponseCode(t, rec.Code, http.StatusOK)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "success" {
			t.Errorf("expected success, got %v", got["status"])
		}
		if got["message"] != "categories fetched successfully" {
			t.Errorf("expected success message, got %v", got["message"])
		}

		// Data should contain 2 categories
		data := got["data"].([]interface{})
		if len(data) != 2 {
			t.Errorf("expected 2 categories, got %d", len(data))
		}
	})

	t.Run("store returns error (500)", func(t *testing.T) {
		store := &tests.StubCategoryStore{Fail: true}
		handler := &categories.Handler{Store: store}

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()

		handler.GetCategories(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "error" {
			t.Errorf("expected error, got %v", got["status"])
		}
	})
}

func createCategoryRequest(data []byte) *http.Request {
	ctx := context.Background()

	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/categories", bytes.NewReader(data))

	return request
}

func assertResponseCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response code = %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want map[string]interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response body = %v, want %v", got, want)
	}
}
