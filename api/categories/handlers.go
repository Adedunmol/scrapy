package categories

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Adedunmol/scrapy/api/helpers"
	"net/http"
)

type Handler struct {
	Store Store
}

func (h *Handler) CreateCategory(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	var body CreateCategoryBody

	// Decode JSON body
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: "error decoding body",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	// Validate request body
	validationErr := helpers.Validate(body)
	if validationErr != nil {
		response := helpers.Response{
			Status:  "error",
			Message: validationErr.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	// Call store method
	category, err := h.Store.CreateCategory(ctx, &body)
	if err != nil {
		// Handle conflict error (e.g., duplicate category name)
		if errors.Is(err, helpers.ErrConflict) {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error(),
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusConflict)
			return
		}

		// Generic internal error
		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}

	// Success response
	response := helpers.Response{
		Status:  "success",
		Message: "category created successfully",
		Data:    category,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
}

func (h *Handler) GetCategories(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	categories, err := h.Store.GetCategories(ctx)
	if err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}

	response := helpers.Response{
		Status:  "success",
		Message: "categories fetched successfully",
		Data:    categories,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
}
