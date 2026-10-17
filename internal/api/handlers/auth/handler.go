package auth

import (
	"context"
	"errors"
	"net/http"

	"icpo-backend/internal/api/handlers"
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
	if err := handlers.DecodeJSONBody(w, r, &input); err != nil {
		return
	}
	if verrs := validator.ValidateRequest(input); verrs != nil {
		handlers.RespondWithValidationErrors(w, verrs)
		return
	}

	user, err := h.service.Register(r.Context(), input.Email, input.Name, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrEmailTaken):
			handlers.RespondWithValidationErrors(w, map[string][]string{
				"email": {err.Error()},
			})
		default:
			handlers.RespondWithServerError(w, err)
		}
		return
	}

	handlers.RespondWithJSON(w, http.StatusCreated, map[string]any{"user": user})
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest
	if err := handlers.DecodeJSONBody(w, r, &input); err != nil {
		return
	}

	if verrs := validator.ValidateRequest(input); verrs != nil {
		handlers.RespondWithValidationErrors(w, verrs)
		return
	}

	pair, err := h.service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			handlers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		default:
			handlers.RespondWithServerError(w, err)
		}
		return
	}

	handlers.SetNoStore(w)
	handlers.RespondWithJSON(w, http.StatusOK, pair)
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input dto.RefreshRequest

	if err := handlers.DecodeJSONBody(w, r, &input); err != nil {
		return
	}

	if input.RefreshToken == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Empty refresh token")
		return
	}

	pair, err := h.service.RefreshToken(r.Context(), input.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidToken):
			handlers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		default:
			handlers.RespondWithServerError(w, err)
		}
		return
	}

	handlers.SetNoStore(w)
	handlers.RespondWithJSON(w, http.StatusOK, pair)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	var input dto.RefreshRequest

	if err := handlers.DecodeJSONBody(w, r, &input); err != nil {
		return
	}

	if input.RefreshToken == "" {
		handlers.RespondWithError(w, http.StatusBadRequest, "Empty refresh token")
		return
	}

	err := h.service.Logout(r.Context(), input.RefreshToken)
	if err != nil {
		handlers.RespondWithServerError(w, err)
		return
	}
	handlers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "logout"})
}
