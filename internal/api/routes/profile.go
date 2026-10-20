package routes

import "net/http"

type ProfileHandler interface {
	Profile(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	UpdatePassword(w http.ResponseWriter, r *http.Request)
}

func registerProfileHandler(mux *http.ServeMux, h ProfileHandler, authenticate Middleware) {
	mux.Handle("GET /profile", authenticate(http.HandlerFunc(h.Profile)))
	mux.Handle("PATCH /profile", authenticate(http.HandlerFunc(h.Update)))
	mux.Handle("PATCH /profile/password", authenticate(http.HandlerFunc(h.UpdatePassword)))
}
