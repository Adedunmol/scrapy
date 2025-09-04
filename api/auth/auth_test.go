package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/tests"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegisterUserHandler(t *testing.T) {
	t.Run("successfully registers a user", func(t *testing.T) {
		store := &tests.StubUserStore{Users: make([]auth.User, 0)}
		handler := &auth.Handler{Store: store}

		data := []byte(`{
			"first_name": "Jane",
			"last_name": "Doe",
			"username": "janedoe",
			"password": "password123",
			"email": "jane@example.com",
			"search_terms": []
		}`)

		req := registerUserRequest(data)
		rec := httptest.NewRecorder()

		handler.RegisterUserHandler(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		want := map[string]interface{}{
			"status":  "success",
			"message": "user created successfully",
		}

		assertResponseCode(t, rec.Code, http.StatusCreated)
		assertResponseBody(t, got, want)

		if len(store.Users) != 1 {
			t.Errorf("expected 1 user in store, got %d", len(store.Users))
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		store := &tests.StubUserStore{}
		handler := &auth.Handler{Store: store}

		// missing closing brace
		data := []byte(`{"first_name": "Jane"`)

		req := registerUserRequest(data)
		rec := httptest.NewRecorder()

		handler.RegisterUserHandler(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)
		assertResponseBody(t, got, map[string]interface{}{
			"status":  "error",
			"message": "error decoding body",
		})
	})

	t.Run("validation error returns 400", func(t *testing.T) {
		store := &tests.StubUserStore{}
		handler := &auth.Handler{Store: store}

		// Missing required email & password
		data := []byte(`{"first_name": "Jane", "last_name": "Doe"}`)

		req := registerUserRequest(data)
		rec := httptest.NewRecorder()

		handler.RegisterUserHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		// message will contain validator error message (stringified)
		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})

	t.Run("store returns conflict (409)", func(t *testing.T) {
		store := &tests.StubUserStore{
			Users: []auth.User{
				{Email: "jane@example.com"}, //ID: 1,
			},
		}
		handler := &auth.Handler{Store: store}

		data := []byte(`{
			"first_name": "Jane",
			"last_name": "Doe",
			"username": "janedoe",
			"password": "password123",
			"email": "jane@example.com",
			"search_terms": []
		}`)

		rec := httptest.NewRecorder()
		req := registerUserRequest(data)
		handler.RegisterUserHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusConflict)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})

	t.Run("store returns generic error (500)", func(t *testing.T) {
		store := &tests.StubUserStore{Fail: true}
		handler := &auth.Handler{Store: store}

		data := []byte(`{
			"first_name": "Jane",
			"last_name": "Doe",
			"username": "janedoe",
			"password": "password123",
			"email": "jane2@example.com",
			"search_terms": []
		}`)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()

		handler.RegisterUserHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})
}

func TestPOSTLogin(t *testing.T) {
	t.Run("returns success when login is valid", func(t *testing.T) {
		store := tests.StubUserStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}
		server := &auth.Handler{Store: &store}

		data := []byte(`{ "email": "adedunmola@gmail.com", "password": "password" }`)
		request := loginUserRequest(data)
		response := httptest.NewRecorder()

		server.LoginUserHandler(response, request)

		var got map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &got)

		assertResponseCode(t, response.Code, http.StatusOK)

		if got["status"] != "success" {
			t.Errorf("expected status success, got %v", got["status"])
		}
		if got["message"] != "user logged in successfully" {
			t.Errorf("expected mocked token, got %v", got["token"])
		}
	})

	t.Run("returns error when body cannot be decoded", func(t *testing.T) {
		store := tests.StubUserStore{Users: []auth.User{}}
		server := &auth.Handler{Store: &store}

		data := []byte(`{ "email": "adedunmola@gmail.com", "password": "password"`) // invalid JSON
		request := loginUserRequest(data)
		response := httptest.NewRecorder()

		server.LoginUserHandler(response, request)

		var got map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &got)

		want := map[string]interface{}{
			"status":  "error",
			"message": "error decoding body",
		}

		assertResponseCode(t, response.Code, http.StatusBadRequest)
		assertResponseBody(t, got, want)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		store := tests.StubUserStore{Users: []auth.User{}}
		server := &auth.Handler{Store: &store}

		data := []byte(`{ "email": "unknown@gmail.com", "password": "password" }`)
		request := loginUserRequest(data)
		response := httptest.NewRecorder()

		server.LoginUserHandler(response, request)

		var got map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &got)

		want := map[string]interface{}{
			"status":  "error",
			"message": "resource not found",
		}

		assertResponseCode(t, response.Code, http.StatusBadRequest)
		assertResponseBody(t, got, want)
	})

	t.Run("returns error when password does not match", func(t *testing.T) {
		store := tests.StubUserStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, // ID: 1,
		}}
		server := &auth.Handler{Store: &store}

		data := []byte(`{ "email": "adedunmola@gmail.com", "password": "wrongpass" }`)
		request := loginUserRequest(data)
		response := httptest.NewRecorder()

		server.LoginUserHandler(response, request)

		var got map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &got)

		want := map[string]interface{}{
			"status":  "error",
			"message": "password does not match",
		}

		assertResponseCode(t, response.Code, http.StatusUnauthorized)
		assertResponseBody(t, got, want)
	})
}

func TestCreateCompanyHandler(t *testing.T) {
	t.Run("successfully creates a company", func(t *testing.T) {
		store := tests.StubUserStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}
		server := &auth.Handler{Store: &store}

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
		store := tests.StubUserStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}
		server := &auth.Handler{Store: &store}

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
		store := tests.StubUserStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		}}
		server := &auth.Handler{Store: &store}

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
		store := tests.StubUserStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		},
			Companies: []auth.Company{
				{Email: "acme@sample.com", Name: "random company"},
			}}
		server := &auth.Handler{Store: &store}

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
		store := tests.StubUserStore{Users: []auth.User{
			{Email: "adedunmola@gmail.com", Password: "password"}, //ID: 1,
		},
			Fail: true}
		server := &auth.Handler{Store: &store}

		data := []byte(`{"name": "Acme Inc", "email": "acme@sample.com"}`)
		req := createCompanyRequest(data, "jane@example.com")
		rec := httptest.NewRecorder()

		server.CreateCompany(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})
}

func registerUserRequest(data []byte) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))

	return req
}

func loginUserRequest(body []byte) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
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
