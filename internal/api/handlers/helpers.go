package handlers

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

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
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

func RespondWithError(w http.ResponseWriter, code int, errorMsg string) {
	if code > 499 {
		slog.Error("Server error", "msg", errorMsg)
	}
	RespondWithJSON(w, code, errorResponse{Message: errorMsg})
}

func RespondWithServerError(w http.ResponseWriter, err error) {
	slog.Error("server error", "err", err)
	RespondWithJSON(w, http.StatusInternalServerError, errorResponse{
		Message: "Internal Server Error",
	})
}

func RespondWithUnauthorizedError(w http.ResponseWriter) {
	RespondWithJSON(w, http.StatusUnauthorized, errorResponse{Message: "Unauthorized"})
}

func RespondWithValidationErrors(w http.ResponseWriter, validationErrs map[string][]string) {
	RespondWithJSON(w, http.StatusUnprocessableEntity, validationErrorsResponse{
		Message: "The given data was invalid",
		Errors:  validationErrs,
	})
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return err
	}
	return nil
}

// SetNoStore запрещает кэширование ответа (RFC 6749 §5.1 — для ответов с токенами).
func SetNoStore(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
}
