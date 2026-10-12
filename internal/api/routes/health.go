package routes

import "net/http"

type HealthHandler interface {
	Check(w http.ResponseWriter, r *http.Request)
}

func registerHealthHandler(mux *http.ServeMux, h HealthHandler) {
	mux.HandleFunc("GET /health", h.Check)
}
