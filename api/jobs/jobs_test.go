package jobs_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Adedunmol/scrapy/api/categories"
	"github.com/Adedunmol/scrapy/api/jobs"
	"github.com/Adedunmol/scrapy/tests"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var UserID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

func TestCreateJobHandler(t *testing.T) {
	t.Run("successfully creates a job", func(t *testing.T) {
		jobStore := &tests.StubJobStore{Companies: []jobs.Company{
			{ID: uuid.New(), Name: "Acme Corp", UserID: UserID},
		}}
		categoryStore := &tests.StubCategoryStore{
			Categories: []categories.Category{
				{ID: uuid.New(), Name: "Engineering"},
			},
		}

		handler := &jobs.Handler{
			Store:           jobStore,
			CategoriesStore: categoryStore,
		}

		data := []byte(`{
			"job_title": "Backend Engineer",
			"job_link": "http://example.com",
			"category": "Engineering",
			"date_posted": "09-10-2025"
		}`)
		req := createJobRequest(data, UserID) // helper to set user_id in context
		rec := httptest.NewRecorder()

		handler.CreateJobHandler(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		assertResponseCode(t, rec.Code, http.StatusCreated)

		if got["status"] != "success" {
			t.Errorf("expected success, got %v", got["status"])
		}

		if len(jobStore.Jobs) != 1 {
			t.Errorf("expected 1 job in store, got %d", len(jobStore.Jobs))
		}
	})

	t.Run("invalid JSON body returns 400", func(t *testing.T) {
		handler := &jobs.Handler{}

		data := []byte(`{"job_title": "Backend Engineer"`) // missing brace
		req := createJobRequest(data, uuid.New())
		rec := httptest.NewRecorder()

		handler.CreateJobHandler(rec, req)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)
		assertResponseBody(t, got, map[string]interface{}{
			"status":  "error",
			"message": "error decoding body",
		})
	})

	t.Run("validation error returns 400", func(t *testing.T) {
		handler := &jobs.Handler{}

		// missing job_title & category
		data := []byte(`{"job_link": "http://example.com"}`)
		req := createJobRequest(data, UserID)
		rec := httptest.NewRecorder()

		handler.CreateJobHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusBadRequest)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		if got["status"] != "error" {
			t.Errorf("expected error status, got %v", got["status"])
		}
	})

	t.Run("category not found returns 404", func(t *testing.T) {
		jobStore := &tests.StubJobStore{
			Companies: []jobs.Company{
				{ID: uuid.New(), Name: "Acme Corp", UserID: UserID},
			},
		}
		categoryStore := &tests.StubCategoryStore{}

		handler := &jobs.Handler{
			Store:           jobStore,
			CategoriesStore: categoryStore,
		}

		data := []byte(`{
			"job_title": "Backend Engineer",
			"job_link": "http://example.com",
			"category": "Engineering",
			"date_posted": "09-10-2025"
		}`)
		req := createJobRequest(data, UserID)
		rec := httptest.NewRecorder()

		handler.CreateJobHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusNotFound)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)

		fmt.Println("got: ")
		fmt.Println(got)

		if got["status"] != "error" {
			t.Errorf("expected error, got %v", got["status"])
		}
	})

	t.Run("company not found returns 404", func(t *testing.T) {
		jobStore := &tests.StubJobStore{NotFound: true}
		categoryStore := &tests.StubCategoryStore{
			Categories: []categories.Category{{ID: uuid.New(), Name: "Engineering"}},
		}

		handler := &jobs.Handler{
			Store:           jobStore,
			CategoriesStore: categoryStore,
		}

		data := []byte(`{
			"job_title": "Backend Engineer",
			"job_link": "http://example.com",
			"category": "Engineering",
			"date_posted": "09-10-2025"
		}`)
		req := createJobRequest(data, uuid.New())
		rec := httptest.NewRecorder()

		handler.CreateJobHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusNotFound)

		var got map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &got)
		if got["status"] != "error" {
			t.Errorf("expected error, got %v", got["status"])
		}
	})

	t.Run("store returns generic error (500)", func(t *testing.T) {
		jobStore := &tests.StubJobStore{Fail: true, Companies: []jobs.Company{
			{ID: uuid.New(), Name: "Acme Corp"},
		}}
		categoryStore := &tests.StubCategoryStore{
			Categories: []categories.Category{{ID: uuid.New(), Name: "Engineering"}},
		}

		handler := &jobs.Handler{
			Store:           jobStore,
			CategoriesStore: categoryStore,
		}

		data := []byte(`{
			"job_title": "Backend Engineer",
			"job_link": "http://example.com",
			"category": "Engineering",
			"date_posted": "09-10-2025"
		}`)
		req := createJobRequest(data, jobStore.Companies[0].ID)
		rec := httptest.NewRecorder()

		handler.CreateJobHandler(rec, req)

		assertResponseCode(t, rec.Code, http.StatusInternalServerError)
	})
}

func createJobRequest(data []byte, userID uuid.UUID) *http.Request {
	ctx := context.WithValue(context.Background(), "user_id", userID)

	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, "/jobs", bytes.NewReader(data))

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
