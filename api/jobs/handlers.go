package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Adedunmol/scrapy/api/categories"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/Adedunmol/scrapy/api/transactions"
	"github.com/Adedunmol/scrapy/api/wallet"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
)

var PerPost = decimal.NewFromInt(100)

type Handler struct {
	Store            Store
	CategoriesStore  categories.Store
	WalletStore      wallet.Store
	TransactionStore transactions.Store
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

	category, err := h.CategoriesStore.GetCategory(ctx, body.Category)
	if err != nil {
		if errors.Is(err, helpers.ErrNotFound) {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error() + ": category not found",
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

	userID := request.Context().Value("user_id")

	company, err := h.Store.GetUserCompany(ctx, userID.(uuid.UUID))
	if err != nil {
		if errors.Is(err, helpers.ErrNotFound) {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error() + ": company not found",
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

	// charge the company's wallet
	walletData, err := h.WalletStore.ChargeWallet(ctx, company.ID, PerPost)
	if err != nil {
		if errors.Is(err, helpers.ErrInsufficientFunds) {
			response := helpers.Response{
				Status:  "error",
				Message: err.Error(),
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
			return
		}
		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}

	// create transaction entry
	txEntry := transactions.CreateTransactionBody{
		Amount:        PerPost,
		BalanceBefore: walletData.Balance.Add(PerPost),
		BalanceAfter:  walletData.Balance,
		Status:        "successful",
		WalletID:      walletData.ID,
		Reference:     uuid.New().String(),
	}
	_, err = h.TransactionStore.CreateTransaction(ctx, &txEntry)

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
		return
	}

	response := helpers.Response{
		Status:  "success",
		Message: "job created successfully",
		Data:    jobData,
	}
	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
	return
}

func (h *Handler) GetUserJobsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	userID := request.Context().Value("user_id")

	q := request.URL.Query()

	pageStr := q.Get("page")
	limitStr := q.Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	jobsData, err := h.Store.GetJobs(ctx, userID.(uuid.UUID), page, limit)
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
		Message: "jobs fetched successfully",
		Data:    jobsData,
	}
	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
	return
}
