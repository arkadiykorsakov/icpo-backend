package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type validationErrorsResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func JSON(w http.ResponseWriter, statusCode int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("marshal response", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		slog.Error("write response", "err", err)
	}
}

func Error(w http.ResponseWriter, code int, msg string) {
	if code >= http.StatusInternalServerError {
		slog.Error("Server error", "msg", msg)
	}
	JSON(w, code, errorResponse{Message: msg})
}

func ServerError(w http.ResponseWriter, err error) {
	slog.Error("server error", "err", err)
	JSON(w, http.StatusInternalServerError, errorResponse{Message: "Internal Server Error"})
}

func Unauthorized(w http.ResponseWriter) {
	JSON(w, http.StatusUnauthorized, errorResponse{Message: "Unauthorized"})
}

func ValidationErrors(w http.ResponseWriter, errs map[string][]string) {
	JSON(w, http.StatusUnprocessableEntity, validationErrorsResponse{
		Message: "The given data was invalid",
		Errors:  errs,
	})
}

// SetNoStore запрещает кэширование ответа (RFC 6749 §5.1 — для ответов с токенами).
func SetNoStore(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
}