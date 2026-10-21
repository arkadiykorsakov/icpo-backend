package routes

import "net/http"

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func registerAuthHandler(mux *http.ServeMux, h AuthHandler, trustedOrigin Middleware) {
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.Handle("POST /auth/refresh", trustedOrigin(http.HandlerFunc(h.Refresh)))
	mux.Handle("POST /auth/logout", trustedOrigin(http.HandlerFunc(h.Logout)))
}
