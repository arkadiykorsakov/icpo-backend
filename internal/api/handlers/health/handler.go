package health

import (
	"net/http"

	"icpo-backend/internal/api/handlers"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Check(w http.ResponseWriter, r *http.Request) {
	handlers.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
