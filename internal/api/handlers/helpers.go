package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("marshal response", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal server error"}`))
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
		slog.Error("server error", "msg", errorMsg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{Error: errorMsg})
}
