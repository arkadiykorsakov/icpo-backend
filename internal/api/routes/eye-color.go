package routes

import "net/http"

type EyeColorHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Autocomplete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func registerEyeColorHandler(mux *http.ServeMux, h EyeColorHandler, authenticate Middleware) {
	mux.Handle("GET /eye-colors", authenticate(http.HandlerFunc(h.GetAll)))
	mux.Handle("GET /eye-colors/autocomplete", authenticate(http.HandlerFunc(h.Autocomplete)))
	mux.Handle("POST /eye-colors", authenticate(http.HandlerFunc(h.Create)))
	mux.Handle("DELETE /eye-colors/{id}", authenticate(http.HandlerFunc(h.Delete)))
}
