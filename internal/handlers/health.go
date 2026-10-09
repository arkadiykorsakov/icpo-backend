package handlers

import "net/http"

func health(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, map[string]string{"status": "ok"})
}
