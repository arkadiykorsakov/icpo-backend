package routes

import "net/http"

type CrimeCategoryHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Autocomplete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func registerCrimeCategoryHandler(mux *http.ServeMux, h CrimeCategoryHandler, authenticate Middleware) {
	mux.Handle("GET /crime-categories", authenticate(http.HandlerFunc(h.GetAll)))
	mux.Handle("GET /crime-categories/autocomplete", authenticate(http.HandlerFunc(h.Autocomplete)))
	mux.Handle("POST /crime-categories", authenticate(http.HandlerFunc(h.Create)))
	mux.Handle("DELETE /crime-categories/{id}", authenticate(http.HandlerFunc(h.Delete)))
}
