package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Adedunmol/scrapy/api/auth"
	"github.com/Adedunmol/scrapy/api/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
)

func CreateSlug(name string) string {
	slug := strings.ToLower(name)

	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}

type Handler struct {
	Store Store
}

func (h *Handler) CreateAdmin(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	var body CreateAdminBody

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

	admin, err := h.Store.CreateAdmin(ctx, &body)
	if err != nil {
		if errors.Is(err, helpers.ErrConflict) {
			response := helpers.Response{
				Status:  "error",
				Message: "admin with this email already exists",
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

	//admin.IsAdmin = true

	response := helpers.Response{
		Status:  "success",
		Message: "admin created successfully",
		Data:    admin,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
}

func (h *Handler) LoginAdminHandler(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	var body auth.LoginUserBody

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

	data, err := h.Store.FindAdminByEmail(ctx, body.Email)
	if err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: err.Error(),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
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

	token, err := helpers.GenerateToken(data.ID, data.Email, true)
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
		Message: "admin logged in successfully",
		Data: map[string]interface{}{
			"token":      token,
			"expiration": helpers.TokenExpiration,
		},
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
}

func (h *Handler) GetAdmins(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	admins, err := h.Store.GetAdmins(ctx)
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
		Message: "admins fetched successfully",
		Data:    admins,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
}

func (h *Handler) CreateRole(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	var body CreateRoleBody

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

	userID := request.Context().Value("user_id").(uuid.UUID)
	body.CreatedBy = userID

	role, err := h.Store.CreateRole(ctx, &body)
	if err != nil {
		if errors.Is(err, helpers.ErrConflict) {
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
		Message: "role created successfully",
		Data:    role,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
}

func (h *Handler) CreatePermission(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	var body CreatePermissionBody

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

	userID, ok := request.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response := helpers.Response{
			Status:  "error",
			Message: "unauthorized: missing user_id",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusUnauthorized)
		return
	}

	body.CreatedBy = userID
	body.Slug = CreateSlug(body.Name)

	permission, err := h.Store.CreatePermission(ctx, &body)
	if err != nil {
		if errors.Is(err, helpers.ErrConflict) {
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
		Message: "permission created successfully",
		Data:    permission,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
}

func (h *Handler) GetRoles(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	roles, err := h.Store.GetRoles(ctx)
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
		Message: "roles fetched successfully",
		Data:    roles,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
}

func (h *Handler) GetPermissions(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	roleID := chi.URLParam(request, "id")
	if roleID == "" {
		response := helpers.Response{
			Status:  "error",
			Message: "role ID is required",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	permissions, err := h.Store.GetPermissions(ctx, roleID)
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
		Message: "permissions fetched successfully",
		Data:    permissions,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusOK)
}

func (h *Handler) BatchCreatePermissions(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	var bodies []CreatePermissionBody

	if err := json.NewDecoder(request.Body).Decode(&bodies); err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: "error decoding body",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	for _, body := range bodies {
		if validationErr := helpers.Validate(body); validationErr != nil {
			response := helpers.Response{
				Status:  "error",
				Message: validationErr.Error(),
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
			return
		}
	}

	userID, ok := request.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response := helpers.Response{
			Status:  "error",
			Message: "unauthorized: missing user_id",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusUnauthorized)
		return
	}

	for i := range bodies {
		bodies[i].CreatedBy = userID
	}

	if err := h.Store.BatchCreatePermissions(ctx, bodies); err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: fmt.Sprintf("failed to batch create permissions: %v", err),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}

	response := helpers.Response{
		Status:  "success",
		Message: "permissions created successfully",
		Data:    nil,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
}

func (h *Handler) BatchCreateRoles(responseWriter http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	var bodies []CreateRoleBody

	if err := json.NewDecoder(request.Body).Decode(&bodies); err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: "error decoding body",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
		return
	}

	for _, body := range bodies {
		if validationErr := helpers.Validate(body); validationErr != nil {
			response := helpers.Response{
				Status:  "error",
				Message: validationErr.Error(),
			}
			helpers.WriteJSONResponse(responseWriter, response, http.StatusBadRequest)
			return
		}
	}

	userID, ok := request.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response := helpers.Response{
			Status:  "error",
			Message: "unauthorized: missing user_id",
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusUnauthorized)
		return
	}

	for i := range bodies {
		bodies[i].CreatedBy = userID
	}

	if err := h.Store.BatchCreateRoles(ctx, bodies); err != nil {
		response := helpers.Response{
			Status:  "error",
			Message: fmt.Sprintf("failed to batch create roles: %v", err),
		}
		helpers.WriteJSONResponse(responseWriter, response, http.StatusInternalServerError)
		return
	}

	response := helpers.Response{
		Status:  "success",
		Message: "roles created successfully",
		Data:    nil,
	}

	helpers.WriteJSONResponse(responseWriter, response, http.StatusCreated)
}
