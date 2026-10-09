package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, errorMsg string) {
	if code > 499 {
		log.Println(errorMsg)
	}

	type errorResponse struct {
		Error string `json:"Error"`
	}

	respondWithJson(w, code, errorResponse{Error: errorMsg})
}
