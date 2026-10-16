package routes

import "net/http"

type TestAuthHandler interface {
	TestAuth(w http.ResponseWriter, r *http.Request)
}

type Authenticate func(http.Handler) http.Handler

func registerTestAuthHandler(mux *http.ServeMux, h TestAuthHandler, authenticate Authenticate) {
	mux.Handle("GET /test-auth", authenticate(http.HandlerFunc(h.TestAuth)))
}
