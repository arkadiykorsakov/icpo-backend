package routes

import "net/http"

type CountryHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Autocomplete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func registerCountryHandler(mux *http.ServeMux, h CountryHandler, authenticate Middleware) {
	mux.Handle("GET /countries", authenticate(http.HandlerFunc(h.GetAll)))
	mux.Handle("GET /countries/autocomplete", authenticate(http.HandlerFunc(h.Autocomplete)))
	mux.Handle("POST /countries", authenticate(http.HandlerFunc(h.Create)))
	mux.Handle("DELETE /countries/{id}", authenticate(http.HandlerFunc(h.Delete)))
}
