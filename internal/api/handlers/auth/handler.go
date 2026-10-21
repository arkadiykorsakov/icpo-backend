package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"icpo-backend/internal/api/handlers/auth/dto"
	"icpo-backend/internal/api/httpx"
	tokenDto "icpo-backend/internal/domain/token/dto"
	userDomain "icpo-backend/internal/domain/user"
	errs "icpo-backend/internal/errors"
	"icpo-backend/pkg/cookies"
	"icpo-backend/pkg/validator"
)

const refreshTokenCookie = "refresh_token"

type Service interface {
	Register(ctx context.Context, email, name, password string) (*userDomain.User, error)
	Login(ctx context.Context, email, password string) (*tokenDto.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*tokenDto.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type handler struct {
	service Service
	cookies cookies.Manager
}

func NewHandler(service Service, cookiesManager cookies.Manager) *handler {
	return &handler{service: service, cookies: cookiesManager}
}

func (h *handler) refreshTokenFromRequest(r *http.Request) (string, error) {
	token, err := h.cookies.Get(r, refreshTokenCookie)
	if errors.Is(err, cookies.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(token), nil
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

	h.cookies.Set(w, refreshTokenCookie, pair.RefreshToken, pair.RefreshExpiresIn)

	httpx.SetNoStore(w)
	httpx.JSON(w, http.StatusOK, pair)
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := h.refreshTokenFromRequest(r)
	if err != nil {
		httpx.ServerError(w, err)
		return
	}
	if refreshToken == "" {
		httpx.Unauthorized(w)
		return
	}

	pair, err := h.service.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidToken) {
			h.cookies.Delete(w, refreshTokenCookie)
			httpx.Unauthorized(w)
			return
		}
		httpx.ServerError(w, err)
		return
	}

	h.cookies.Set(w, refreshTokenCookie, pair.RefreshToken, pair.RefreshExpiresIn)

	httpx.SetNoStore(w)
	httpx.JSON(w, http.StatusOK, pair)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := h.refreshTokenFromRequest(r)
	if err != nil {
		httpx.ServerError(w, err)
		return
	}

	if refreshToken != "" {
		if err := h.service.Logout(r.Context(), refreshToken); err != nil {
			httpx.ServerError(w, err)
			return
		}
	}

	h.cookies.Delete(w, refreshTokenCookie)
	w.WriteHeader(http.StatusNoContent)
}
