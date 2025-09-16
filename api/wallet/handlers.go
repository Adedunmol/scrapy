package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/http"
)

type Handler struct {
	Store Store
}

func (h *Handler) GetWalletHandler(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	userID := request.Context().Value("user_id")

	wallet, err := h.Store.GetWallet(ctx, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, helpers.ErrNotFound) {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error(),
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusNotFound)
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
		Message: "retrieved company's wallet successfully",
		Data:    wallet,
	}
	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
	return
}

func (h *Handler) TopUpWalletHandler(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	var body TopUpWalletBody

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

	if body.Amount.LessThanOrEqual(decimal.Zero) {
		response := helpers.Response{
			Status:  "error",
			Message: "amount must be greater than zero",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	userID := request.Context().Value("user_id")

	wallet, err := h.Store.TopUpWallet(ctx, userID.(uuid.UUID), body.Amount)

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
		Message: "updated company's wallet successfully",
		Data:    wallet,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
	return
}
