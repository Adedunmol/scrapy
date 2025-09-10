package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Adedunmol/scrapy/api/categories"
	"github.com/Adedunmol/scrapy/api/companies"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	Store           Store
	categoriesStore categories.Store
	companiesStore  companies.Store
}

func (h *Handler) CreateJobHandler(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	var body CreateJobBody

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: "error decoding body",
		}

		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	validationErr := helpers.Validate(body)
	if validationErr != nil {
		response := helpers.Response{
			Status:  "error",
			Message: validationErr.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	category, err := h.categoriesStore.GetCategory(ctx, body.Category)
	if err != nil {
		if errors.Is(err, helpers.ErrNotFound) {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error() + ": category not found",
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusNotFound)
		}
		response := helpers.Response{
			Status:  "error",
			Message: http.StatusText(http.StatusInternalServerError),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
	}

	userID := request.Context().Value("user_id")

	company, err := h.companiesStore.GetUserCompany(ctx, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, helpers.ErrNotFound) {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error() + ": company not found",
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusNotFound)
		}
		response := helpers.Response{
			Status:  "error",
			Message: http.StatusText(http.StatusInternalServerError),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
	}

	body.Origin = "company"
	body.OriginID = company.ID
	body.CategoryID = category.ID

	jobData, err := h.Store.CreateJob(ctx, &body)
	if err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
	}

	response := helpers.Response{
		Status:  "success",
		Message: "job created successfully",
		Data:    jobData,
	}
	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
}
