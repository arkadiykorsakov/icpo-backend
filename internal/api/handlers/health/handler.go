package health

import (
	"net/http"

	"icpo-backend/internal/api/httpx"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Check(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
