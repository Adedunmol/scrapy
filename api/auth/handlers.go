package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/Adedunmol/scrapy/queue"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Handler struct {
	Store Store
	Queue queue.Queue
}

/*
RegisterUserHandler is to create the user and check the prefilled categories table
for valid search terms the user picked and then create entries in the preferences table for the user
*/
func (h *Handler) RegisterUserHandler(responseWriter http.ResponseWriter, request *http.Request) {

	ctx := context.Background()
	var body CreateUserBody

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}

	body.Password = string(hashedPassword)

	categories, err := h.Store.GetCategories(ctx)
	if err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}

	var validSearchTerms []uuid.UUID

	for _, uc := range body.SearchTerms {
		if id, ok := categories[uc]; ok {
			validSearchTerms = append(validSearchTerms, id)
		}
	}

	user, err := h.Store.CreateUser(&body)

	if err != nil {
		ok := errors.As(err, &helpers.ErrConflict)

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

	err = h.Store.CreatePreferences(ctx, validSearchTerms, user.ID)
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
		Message: "user created successfully",
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
	return
}

func (h *Handler) LoginUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var body LoginUserBody

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

	data, err := h.Store.FindUserByEmail(body.Email)
	if err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusNotFound)
		return
	}

	match := h.Store.ComparePasswords(data.Password, body.Password)
	if !match {
		response := helpers.Response{
			Status:  "error",
			Message: "password does not match",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusUnauthorized)
		return
	}
	token, err := helpers.GenerateToken(data.ID, data.Email)
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
		Message: "user logged in successfully",
		Data:    map[string]interface{}{"token": token, "expiration": helpers.TokenExpiration},
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
}
