package handlers

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)

	mux.HandleFunc("GET /swagger", swaggerRedirect)
	mux.HandleFunc("GET /swagger/{$}", swaggerUIHandler)
	mux.HandleFunc("GET /swagger/openapi.yaml", openAPIHandler)

	return loggingMiddleware(mux)
}
