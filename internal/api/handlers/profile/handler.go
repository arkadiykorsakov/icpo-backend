package profile

import (
	"context"
	"errors"
	"net/http"

	"icpo-backend/internal/api/authctx"
	"icpo-backend/internal/api/handlers/profile/dto"
	"icpo-backend/internal/api/httpx"
	userDomain "icpo-backend/internal/domain/user"
	errs "icpo-backend/internal/errors"
	"icpo-backend/pkg/validator"
)

type Service interface {
	Profile(ctx context.Context, userID int64) (userDomain.User, error)
	Update(ctx context.Context, userID int64, email, name string) (userDomain.User, error)
	UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := authctx.UserID(r.Context())
	if !ok {
		httpx.ServerError(w, errs.ErrUserIDMissing)
		return
	}

	user, err := h.service.Profile(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			httpx.Error(w, http.StatusNotFound, "User not found")
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"profile": user,
	})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := authctx.UserID(r.Context())
	if !ok {
		httpx.ServerError(w, errs.ErrUserIDMissing)
		return
	}

	var input dto.UpdateProfileRequest
	if err := httpx.DecodeJSON(w, r, &input); err != nil {
		return
	}

	if verrs := validator.ValidateRequest(input); verrs != nil {
		httpx.ValidationErrors(w, verrs)
		return
	}

	userUpdated, err := h.service.Update(r.Context(), userID, input.Email, input.Name)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			httpx.Error(w, http.StatusNotFound, "User not found")
		case errors.Is(err, errs.ErrEmailTaken):
			httpx.ValidationErrors(w, map[string][]string{
				"email": {err.Error()},
			})
		default:
			httpx.ServerError(w, err)
		}
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{
		"profile": userUpdated,
	})
}

func (h *handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := authctx.UserID(r.Context())
	if !ok {
		httpx.ServerError(w, errs.ErrUserIDMissing)
		return
	}

	var input dto.UpdatePasswordRequest
	if err := httpx.DecodeJSON(w, r, &input); err != nil {
		return
	}

	if verrs := validator.ValidateRequest(input); verrs != nil {
		httpx.ValidationErrors(w, verrs)
		return
	}

	err := h.service.UpdatePassword(r.Context(), userID, input.OldPassword, input.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			httpx.Error(w, http.StatusNotFound, "User not found")
		case errors.Is(err, errs.ErrWrongPassword):
			httpx.ValidationErrors(w, map[string][]string{
				"old_password": {err.Error()},
			})
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
