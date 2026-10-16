package testauth

import (
	"net/http"

	"icpo-backend/internal/api/authctx"
	"icpo-backend/internal/api/handlers"
	errs "icpo-backend/internal/errors"
)

type handler struct{}

func NewHandler() *handler { return &handler{} }

func (h *handler) TestAuth(w http.ResponseWriter, r *http.Request) {
	userID, ok := authctx.UserID(r.Context())
	if !ok {
		handlers.RespondWithServerError(w, errs.ErrUserIDMissing)
		return
	}

	handlers.RespondWithJSON(w, http.StatusOK, map[string]any{
		"user_id": userID,
	})
}
