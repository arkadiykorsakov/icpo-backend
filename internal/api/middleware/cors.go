package middleware

import "net/http"

func originSet(origins []string) map[string]struct{} {
	set := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		set[o] = struct{}{}
	}
	return set
}

// CORS разрешает запросы с cookie только с origin из белого списка.
// Должен стоять снаружи mux, чтобы отвечать на preflight до маршрутизации.
func CORS(allowed []string) func(http.Handler) http.Handler {
	set := originSet(allowed)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Add("Vary", "Origin")

			origin := r.Header.Get("Origin")
			if _, ok := set[origin]; ok {
				h.Set("Access-Control-Allow-Origin", origin)
				h.Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				if _, ok := set[origin]; ok {
					h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
					h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
					h.Set("Access-Control-Max-Age", "600")
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
