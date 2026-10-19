package testauth

import (
	"net/http"

	"icpo-backend/internal/api/authctx"
	"icpo-backend/internal/api/httpx"
	errs "icpo-backend/internal/errors"
)

type handler struct{}

func NewHandler() *handler { return &handler{} }

func (h *handler) TestAuth(w http.ResponseWriter, r *http.Request) {
	userID, ok := authctx.UserID(r.Context())
	if !ok {
		httpx.ServerError(w, errs.ErrUserIDMissing)
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"user_id": userID,
	})
}
