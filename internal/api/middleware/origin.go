package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"icpo-backend/internal/api/httpx"
)

// RequireTrustedOrigin защищает cookie-ручки от CSRF. Если браузер прислал Origin,
// он должен быть либо в белом списке, либо совпадать с хостом самого сервиса
// (same-origin: фронт за reverse proxy, Swagger UI). Запросы без Origin
// (curl, мобильные клиенты) пропускаются: CSRF возможен только из браузера.
func RequireTrustedOrigin(allowed []string) func(http.Handler) http.Handler {
	set := originSet(allowed)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get("Origin"); origin != "" && !isTrustedOrigin(origin, r.Host, set) {
				httpx.Error(w, http.StatusForbidden, "Forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isTrustedOrigin(origin, host string, allowed map[string]struct{}) bool {
	if _, ok := allowed[origin]; ok {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return false // в том числе "Origin: null"
	}
	return strings.EqualFold(u.Host, host)
}
