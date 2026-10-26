package routes

import "net/http"

type LanguageHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Autocomplete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func registerLanguageHandler(mux *http.ServeMux, h LanguageHandler, authenticate Middleware) {
	mux.Handle("GET /languages", authenticate(http.HandlerFunc(h.GetAll)))
	mux.Handle("GET /languages/autocomplete", authenticate(http.HandlerFunc(h.Autocomplete)))
	mux.Handle("POST /languages", authenticate(http.HandlerFunc(h.Create)))
	mux.Handle("DELETE /languages/{id}", authenticate(http.HandlerFunc(h.Delete)))
}
