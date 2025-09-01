package auth

import (
	"encoding/json"
	"errors"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/Adedunmol/scrapy/queue"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Handler struct {
	Store Store
	Queue queue.Queue
}

const OtpExpiration = 30

func (h *Handler) CreateUserHandler(responseWriter http.ResponseWriter, request *http.Request) {

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
	err = h.Store.CreateUser(body)

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

	response := helpers.Response{
		Status:  "success",
		Message: "user created successfully",
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
	return
}

//func (h *Handler) LoginUserHandler(responseWriter http.ResponseWriter, request *http.Request) {
//	body, problems, err := helpers.DecodeAndValidate[*LoginUserBody](request)
//
//	var clientError helpers.ClientError
//	ok := errors.As(err, &clientError)
//
//	if err != nil && problems == nil {
//		helpers.HandleError(responseWriter, helpers.NewHTTPError(err, http.StatusBadRequest, "invalid request body", nil))
//		return
//	}
//
//	if err != nil && ok {
//		helpers.HandleError(responseWriter, helpers.NewHTTPError(err, http.StatusBadRequest, "invalid request body", problems))
//		return
//	}
//
//	data, err := h.Store.FindUserByEmail(body.Email)
//	if err != nil {
//		helpers.HandleError(responseWriter, helpers.ErrUnauthorized)
//		return
//	}
//
//	match := h.Store.ComparePasswords(data.Password, body.Password)
//	if !match {
//		helpers.HandleError(responseWriter, helpers.ErrUnauthorized)
//		return
//	}
//
//	token, err := helpers.GenerateToken(data.ID, data.Email, data.Verified)
//	if err != nil {
//		helpers.HandleError(responseWriter, helpers.NewHTTPError(err, http.StatusInternalServerError, "internal server error", nil))
//		return
//	}
//
//	refreshToken, err := helpers.GenerateToken(data.ID, data.Email, data.Verified)
//	if err != nil {
//		helpers.HandleError(responseWriter, helpers.NewHTTPError(err, http.StatusInternalServerError, "internal server error", nil))
//		return
//	}
//
//	updateUser := UpdateUserBody{RefreshToken: refreshToken}
//
//	if _, err = h.Store.UpdateUser(data.ID, updateUser); err != nil {
//		helpers.HandleError(responseWriter, helpers.ErrInternalServerError)
//		return
//	}
//
//	expires := time.Now().AddDate(0, 1, 0)
//
//	cookie := &http.Cookie{
//		Name:     "refresh_token",
//		Value:    refreshToken,
//		Path:     "/",
//		Expires:  expires,
//		Secure:   true,
//		HttpOnly: true,
//		MaxAge:   86400,
//	}
//
//	http.SetCookie(responseWriter, cookie)
//
//	response := Response{
//		Status:  "Success",
//		Message: "User logged in",
//		Data:    map[string]interface{}{"token": token, "expiration": helpers.TokenExpiration},
//	}
//
//	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
//}
