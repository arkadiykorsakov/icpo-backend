package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"icpo-backend/internal/api/authctx"
	"icpo-backend/internal/api/httpx"
	tokenDomain "icpo-backend/internal/domain/token"
)

type TokenService interface {
	ParseAccessToken(tokenStr string) (*tokenDomain.Claims, error)
}

func Authenticate(tokenSvc TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := extractBearerToken(r.Header.Get("Authorization"))
			if !ok {
				slog.Debug("missing or malformed Authorization header",
					"method", r.Method, "path", r.URL.Path)
				unauthorized(w)
				return
			}

			claims, err := tokenSvc.ParseAccessToken(tokenString)
			if err != nil || claims == nil {
				slog.Debug("invalid access token",
					"err", err, "method", r.Method, "path", r.URL.Path)
				unauthorized(w)
				return
			}

			ctx := authctx.WithUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(authorization string) (string, bool) {
	scheme, token, ok := strings.Cut(authorization, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}
	return token, true
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
	httpx.Unauthorized(w)
}
