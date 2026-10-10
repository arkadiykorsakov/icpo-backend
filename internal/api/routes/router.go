package routes

import (
	"net/http"

	"icpo-backend/internal/api/handlers"
	"icpo-backend/internal/api/middleware"
)

func NewRouter() http.Handler {
	api := http.NewServeMux()
	api.HandleFunc("GET /health", handlers.Health)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))

	mux.HandleFunc("GET /swagger", handlers.SwaggerRedirect)
	mux.HandleFunc("GET /swagger/{$}", handlers.SwaggerUI)
	mux.HandleFunc("GET /swagger/openapi.yaml", handlers.OpenAPI)

	return middleware.Logging(mux)
}
