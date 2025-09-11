package companies

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Handler struct {
	Store Store
}

func (h *Handler) CreateCompany(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	userID := request.Context().Value("user_id")

	var body CreateCompanyBody

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

	body.UserID = userID.(uuid.UUID)

	company, err := h.Store.CreateCompany(ctx, &body)

	if err != nil {
		ok := errors.Is(err, helpers.ErrConflict)

		log.Println("ok: ", ok)
		if ok {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error(),
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusConflict)
			return
		}

		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}

	response := helpers.Response{
		Status:  "success",
		Message: "company created successfully",
		Data:    company,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
	return
}
