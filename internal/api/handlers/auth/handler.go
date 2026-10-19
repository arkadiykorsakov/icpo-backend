package auth

import (
	"context"
	"errors"
	"net/http"

	"icpo-backend/internal/api/httpx"
	"icpo-backend/internal/api/handlers/auth/dto"
	tokenDto "icpo-backend/internal/domain/token/dto"
	userDomain "icpo-backend/internal/domain/user"
	errs "icpo-backend/internal/errors"
	"icpo-backend/pkg/validator"
)

type Service interface {
	Register(ctx context.Context, email, name, password string) (*userDomain.User, error)
	Login(ctx context.Context, email, password string) (*tokenDto.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*tokenDto.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterRequest
	if err := httpx.DecodeJSON(w, r, &input); err != nil {
		return
	}
	if verrs := validator.ValidateRequest(input); verrs != nil {
		httpx.ValidationErrors(w, verrs)
		return
	}

	user, err := h.service.Register(r.Context(), input.Email, input.Name, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrEmailTaken):
			httpx.ValidationErrors(w, map[string][]string{
				"email": {err.Error()},
			})
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, map[string]any{"user": user})
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest
	if err := httpx.DecodeJSON(w, r, &input); err != nil {
		return
	}

	if verrs := validator.ValidateRequest(input); verrs != nil {
		httpx.ValidationErrors(w, verrs)
		return
	}

	pair, err := h.service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			httpx.Error(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	httpx.SetNoStore(w)
	httpx.JSON(w, http.StatusOK, pair)
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input dto.RefreshRequest

	if err := httpx.DecodeJSON(w, r, &input); err != nil {
		return
	}

	if input.RefreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "Empty refresh token")
		return
	}

	pair, err := h.service.RefreshToken(r.Context(), input.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidToken):
			httpx.Error(w, http.StatusUnauthorized, err.Error())
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	httpx.SetNoStore(w)
	httpx.JSON(w, http.StatusOK, pair)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	var input dto.RefreshRequest

	if err := httpx.DecodeJSON(w, r, &input); err != nil {
		return
	}

	if input.RefreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "Empty refresh token")
		return
	}

	err := h.service.Logout(r.Context(), input.RefreshToken)
	if err != nil {
		httpx.ServerError(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"message": "logout"})
}
