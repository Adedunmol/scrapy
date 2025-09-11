package companies_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/companies"
	"github.com/Adedunmol/scrapy/tests"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateCompanyHandler(t *testing.T) {
	t.Run("successfully creates a company", func(t *testing.T) {
		store := tests.StubCompanyStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}
		server := &companies.Handler{Store: &store}

		data := []byte(`{"name": "Acme Inc", "email": "acme@sample.com"}`)
		req := createCompanyRequest(data, "jane@example.com")
		rec := httptest.NewRecorder()

		server.CreateCompany(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		assertResponseCode(t, rec.Code, http.StatusCreated)

		if len(store.Companies) != 1 {
			t.Errorf("expected 1 company in store, got %d", len(store.Companies))
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		store := tests.StubCompanyStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}
		server := &companies.Handler{Store: &store}

		// missing closing brace
		data := []byte(`{"name": "Acme Inc"`)
		req := createCompanyRequest(data, "jane@example.com")
		rec := httptest.NewRecorder()

		server.CreateCompany(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)
		assertResponseBody(t, got, map[string]interface{}{
			"status":  "error",
			"message": "error decoding body",
		})
	})

	t.Run("validation error returns 400", func(t *testing.T) {
		store := tests.StubCompanyStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}
		server := &companies.Handler{Store: &store}

		// Missing required fields
		data := []byte(`{}`)
		req := createCompanyRequest(data, "jane@example.com")
		rec := httptest.NewRecorder()

		server.CreateCompany(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)
		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})

	t.Run("store returns conflict (409)", func(t *testing.T) {
		store := tests.StubCompanyStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		},
			Companies: []companies.Company{
				{Email: "acme@sample.com", Name: "random company"},
			}}
		server := &companies.Handler{Store: &store}

		data := []byte(`{"name": "Acme Inc", "email": "acme@sample.com"}`)
		req := createCompanyRequest(data, "jane@example.com")
		rec := httptest.NewRecorder()

		server.CreateCompany(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		assertResponseCode(t, rec.Code, http.StatusConflict)
		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})

	t.Run("store returns generic error (500)", func(t *testing.T) {
		store := tests.StubCompanyStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		},
			Fail: true}
		server := &companies.Handler{Store: &store}

		data := []byte(`{"name": "Acme Inc", "email": "acme@sample.com"}`)
		req := createCompanyRequest(data, "jane@example.com")
		rec := httptest.NewRecorder()

		server.CreateCompany(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})
}

func createCompanyRequest(data []byte, email string) *http.Request {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), "email", email)
	ctx = context.WithValue(context.Background(), "user_id", userID)

	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/companies", bytes.NewReader(data))

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
