package routes

import (
	"net/http"

	"icpo-backend/internal/api/middleware"
)

type Middleware = func(http.Handler) http.Handler

type Deps struct {
	Health        HealthHandler
	Swagger       SwaggerHandler
	Auth          AuthHandler
	Profile       ProfileHandler
	Country       CountryHandler
	Language      LanguageHandler
	CrimeCategory CrimeCategoryHandler
	HairColor     HairColorHandler
	EyeColor      EyeColorHandler
	City          CityHandler
	Authenticate  Middleware
	CORS          Middleware
	TrustedOrigin Middleware
}

func NewRouter(d Deps) http.Handler {
	api := http.NewServeMux()
	registerHealthHandler(api, d.Health)
	registerAuthHandler(api, d.Auth, d.TrustedOrigin)
	registerProfileHandler(api, d.Profile, d.Authenticate)
	registerCountryHandler(api, d.Country, d.Authenticate)
	registerLanguageHandler(api, d.Language, d.Authenticate)
	registerCrimeCategoryHandler(api, d.CrimeCategory, d.Authenticate)
	registerHairColorHandler(api, d.HairColor, d.Authenticate)
	registerEyeColorHandler(api, d.EyeColor, d.Authenticate)
	registerCityHandler(api, d.City, d.Authenticate)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))
	registerSwaggerHandler(mux, d.Swagger)

	return middleware.Logging(d.CORS(mux))
}
