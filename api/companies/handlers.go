package companies

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/Adedunmol/scrapy/api/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Handler struct {
	Store       Store
	WalletStore wallet.Store
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

	// create company wallet
	_, err = h.WalletStore.CreateWallet(ctx, company.ID)
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
		Message: "company created successfully",
		Data:    company,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
	return
}

func (h *Handler) GetCompany(responseWriter http.ResponseWriter, request *http.Request) {
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

	company, err := h.Store.GetCompany(ctx, uuid.MustParse(companyID))
	if err != nil {
		log.Println(err)
		ok := errors.Is(err, helpers.ErrNotFound)
		if ok {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error(),
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusNotFound)
			return
		}
		response := helpers.Response{
			Status:  "error",
			Message: http.StatusText(http.StatusInternalServerError),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}
	response := helpers.Response{
		Status:  "success",
		Message: "company retrieved successfully",
		Data:    company,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
	return
}

func (h *Handler) GetCompanyJobsHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	companiesData, err := h.Store.GetCompanyJobs(ctx, uuid.MustParse(companyID))
	if err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: http.StatusText(http.StatusInternalServerError),
		}

		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}
	response := helpers.Response{
		Status:  "success",
		Message: "company jobs retrieved successfully",
		Data:    companiesData,
	}
	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
	return
}
