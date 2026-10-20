package routes

import (
	"net/http"

	"icpo-backend/internal/api/middleware"
)

type Middleware = func(http.Handler) http.Handler

type Deps struct {
	Health       HealthHandler
	Swagger      SwaggerHandler
	Auth         AuthHandler
	Profile      ProfileHandler
	Authenticate Middleware
}

func NewRouter(d Deps) http.Handler {
	api := http.NewServeMux()
	registerHealthHandler(api, d.Health)
	registerAuthHandler(api, d.Auth)
	registerProfileHandler(api, d.Profile, d.Authenticate)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))
	registerSwaggerHandler(mux, d.Swagger)

	return middleware.Logging(mux)
}
