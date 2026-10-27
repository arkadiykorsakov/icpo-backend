package routes

import "net/http"

type HairColorHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Autocomplete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func registerHairColorHandler(mux *http.ServeMux, h HairColorHandler, authenticate Middleware) {
	mux.Handle("GET /hair-colors", authenticate(http.HandlerFunc(h.GetAll)))
	mux.Handle("GET /hair-colors/autocomplete", authenticate(http.HandlerFunc(h.Autocomplete)))
	mux.Handle("POST /hair-colors", authenticate(http.HandlerFunc(h.Create)))
	mux.Handle("DELETE /hair-colors/{id}", authenticate(http.HandlerFunc(h.Delete)))
}
