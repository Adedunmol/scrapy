package wallet

import (
	"context"
	"errors"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	Store Store
}

func (h *Handler) GetWalletHandler(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	companyID := chi.URLParam(request, "company_id")

	if companyID == "" {
		response := helpers.Response{
			Status:  "error",
			Message: "company_id is required",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	wallet, err := h.Store.GetWallet(ctx, uuid.MustParse(companyID))
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
