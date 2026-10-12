package routes

import (
	"net/http"

	"icpo-backend/internal/api/middleware"
)

type Deps struct {
	Health  HealthHandler
	Swagger SwaggerHandler
}

func NewRouter(d Deps) http.Handler {
	api := http.NewServeMux()
	registerHealthHandler(api, d.Health)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))
	registerSwaggerHandler(mux, d.Swagger)

	return middleware.Logging(mux)
}
