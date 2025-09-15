package companies_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/companies"
	"github.com/Adedunmol/scrapy/api/jobs"
	"github.com/Adedunmol/scrapy/tests"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateCompanyHandler(t *testing.T) {
	t.Run("successfully creates a company", func(t *testing.T) {
		id := uuid.New()
		store := tests.StubCompanyStore{Users: []auth.User{
			{ID: id, Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}
		walletStore := tests.StubWalletStore{}
		server := &companies.Handler{Store: &store, WalletStore: &walletStore}

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
		id := uuid.New()

		store := tests.StubCompanyStore{Users: []auth.User{
			{ID: id, Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}

		walletStore := tests.StubWalletStore{}
		server := &companies.Handler{Store: &store, WalletStore: &walletStore}

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
		id := uuid.New()

		store := tests.StubCompanyStore{Users: []auth.User{
			{ID: id, Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}

		walletStore := tests.StubWalletStore{}
		server := &companies.Handler{Store: &store, WalletStore: &walletStore}

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
		id := uuid.New()

		store := tests.StubCompanyStore{Users: []auth.User{
			{ID: id, Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		},
			Companies: []companies.Company{
				{Email: "acme@sample.com", Name: "random company"},
			}}
		walletStore := tests.StubWalletStore{}
		server := &companies.Handler{Store: &store, WalletStore: &walletStore}

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
		id := uuid.New()

		store := tests.StubCompanyStore{Users: []auth.User{
			{ID: id, Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		},
			Fail: true}
		walletStore := tests.StubWalletStore{}
		server := &companies.Handler{Store: &store, WalletStore: &walletStore}

		data := []byte(`{"name": "Acme Inc", "email": "acme@sample.com"}`)
		req := createCompanyRequest(data, "jane@example.com")
		rec := httptest.NewRecorder()

		server.CreateCompany(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})

	t.Run("store returns generic error (500) for wallet", func(t *testing.T) {
		id := uuid.New()

		store := tests.StubCompanyStore{Users: []auth.User{
			{ID: id, Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		},
			Fail: true}
		walletStore := tests.StubWalletStore{}
		server := &companies.Handler{Store: &store, WalletStore: &walletStore}

		data := []byte(`{"name": "Acme Inc", "email": "acme@sample.com"}`)
		req := createCompanyRequest(data, "jane@example.com")
		rec := httptest.NewRecorder()

		server.CreateCompany(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestGetCompanyHandler(t *testing.T) {
	t.Run("successfully retrieves a company", func(t *testing.T) {
		companyID := uuid.New()
		store := &tests.StubCompanyStore{
			Companies: []companies.Company{
				{ID: companyID, Name: "Acme Inc", Email: "acme@sample.com"},
			},
		}
		handler := &companies.Handler{Store: store}

		req := getCompanyRequest(companyID)
		rec := httptest.NewRecorder()

		handler.GetCompany(rec, req)

		assertResponseCode(t, rec.Code, http.StatusOK)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "success" {
			t.Errorf("expected status success, got %v", got["status"])
		}
		if got["message"] != "company retrieved successfully" {
			t.Errorf("expected message company retrieved successfully, got %v", got["message"])
		}
	})

	t.Run("missing company_id returns 400", func(t *testing.T) {
		store := &tests.StubCompanyStore{}
		handler := &companies.Handler{Store: store}

		req := httptest.NewRequest(http.MethodGet, "/companies/", nil)
		rec := httptest.NewRecorder()

		// no param injected here
		handler.GetCompany(rec, req)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["message"] != "company_id is required" {
			t.Errorf("expected company_id is required, got %v", got["message"])
		}
	})

	t.Run("company not found returns 404", func(t *testing.T) {
		companyID := uuid.New()
		store := &tests.StubCompanyStore{NotFound: true}
		handler := &companies.Handler{Store: store}

		req := getCompanyRequest(companyID)
		rec := httptest.NewRecorder()

		handler.GetCompany(rec, req)

		assertResponseCode(t, rec.Code, http.StatusNotFound)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})

	t.Run("store failure returns 500", func(t *testing.T) {
		companyID := uuid.New()
		store := &tests.StubCompanyStore{Fail: true}
		handler := &companies.Handler{Store: store}

		req := getCompanyRequest(companyID)

		rec := httptest.NewRecorder()

		handler.GetCompany(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestGetCompanyJobsHandler(t *testing.T) {
	t.Run("successfully retrieves company jobs", func(t *testing.T) {
		companyID := uuid.New()
		store := &tests.StubCompanyStore{
			Jobs: []jobs.Job{
				{ID: uuid.New(), JobTitle: "Software Engineer", OriginID: companyID},
				{ID: uuid.New(), JobTitle: "Product Manager", OriginID: companyID},
			},
			Companies: []companies.Company{
				{ID: companyID, Name: "Acme Inc", Email: "acme@sample.com"},
			},
		}
		handler := &companies.Handler{Store: store}

		req := getCompanyJobsRequest(companyID)
		rec := httptest.NewRecorder()

		handler.GetCompanyJobsHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusOK)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "success" {
			t.Errorf("expected status success, got %v", got["status"])
		}
		if got["message"] != "company jobs retrieved successfully" {
			t.Errorf("expected message company jobs retrieved successfully, got %v", got["message"])
		}
		data := got["data"].([]interface{})
		if len(data) != 2 {
			t.Errorf("expected 2 jobs, got %d", len(data))
		}
	})

	t.Run("missing company_id returns 400", func(t *testing.T) {
		store := &tests.StubCompanyStore{}
		handler := &companies.Handler{Store: store}

		req := httptest.NewRequest(http.MethodGet, "/companies//jobs", nil)
		rec := httptest.NewRecorder()

		// no company_id injected
		handler.GetCompanyJobsHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["message"] != "company_id is required" {
			t.Errorf("expected company_id is required, got %v", got["message"])
		}
	})

	t.Run("store failure returns 500", func(t *testing.T) {
		companyID := uuid.New()
		store := &tests.StubCompanyStore{Fail: true}
		handler := &companies.Handler{Store: store}

		req := getCompanyJobsRequest(companyID)
		rec := httptest.NewRecorder()

		handler.GetCompanyJobsHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})
}

func getCompanyRequest(companyID uuid.UUID) *http.Request {
	ctx := context.Background()

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/companies/"+companyID.String(), nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("company_id", companyID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req
}
func getCompanyJobsRequest(companyID uuid.UUID) *http.Request {
	ctx := context.Background()

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/companies/"+companyID.String()+"/jobs", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("company_id", companyID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req
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
