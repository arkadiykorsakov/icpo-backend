package routes

import "net/http"

type CityHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Autocomplete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func registerCityHandler(mux *http.ServeMux, h CityHandler, authenticate Middleware) {
	mux.Handle("GET /cities", authenticate(http.HandlerFunc(h.GetAll)))
	mux.Handle("GET /cities/autocomplete", authenticate(http.HandlerFunc(h.Autocomplete)))
	mux.Handle("POST /cities", authenticate(http.HandlerFunc(h.Create)))
	mux.Handle("DELETE /cities/{id}", authenticate(http.HandlerFunc(h.Delete)))
}
