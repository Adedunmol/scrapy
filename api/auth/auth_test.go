package auth_test

import (
	"bytes"
	"encoding/json"
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/tests"
	"net/http"
	"net/http/httptest"
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
			"email": "jane@example.com"
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

		tests.AssertResponseCode(t, rec.Code, http.StatusCreated)
		tests.AssertResponseBody(t, got, want)

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

		tests.AssertResponseCode(t, rec.Code, http.StatusBadRequest)
		tests.AssertResponseBody(t, got, map[string]interface{}{
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

		tests.AssertResponseCode(t, rec.Code, http.StatusBadRequest)

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
				{ID: 1, Email: "jane@example.com"},
			},
		}
		handler := &auth.Handler{Store: store}

		data := []byte(`{
			"first_name": "Jane",
			"last_name": "Doe",
			"username": "janedoe",
			"password": "password123",
			"email": "jane@example.com"
		}`)

		rec := httptest.NewRecorder()
		req := registerUserRequest(data)
		handler.RegisterUserHandler(rec, req)

		tests.AssertResponseCode(t, rec.Code, http.StatusConflict)

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
			"email": "jane2@example.com"
		}`)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))
		rec := httptest.NewRecorder()

		handler.RegisterUserHandler(rec, req)

		tests.AssertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})
}

func registerUserRequest(data []byte) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))

	return req
}
