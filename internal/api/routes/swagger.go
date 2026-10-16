package routes

import "net/http"

type SwaggerHandler interface {
	Redirect(w http.ResponseWriter, r *http.Request)
	UI(w http.ResponseWriter, r *http.Request)
	OpenAPI(w http.ResponseWriter, r *http.Request)
}

func registerSwaggerHandler(mux *http.ServeMux, h SwaggerHandler) {
	mux.HandleFunc("GET /swagger", h.Redirect)
	mux.HandleFunc("GET /swagger/{$}", h.UI)
	mux.HandleFunc("GET /swagger/openapi.yaml", h.OpenAPI)
}
