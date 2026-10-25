package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func GetIDFromRequest(w http.ResponseWriter, r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		Error(w, http.StatusBadRequest, "Id is required")
		return 0, errors.New("id is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id < 1 {
		Error(w, http.StatusBadRequest, "Id must be a positive integer")
		return 0, errors.New("id must be a positive integer")
	}

	return id, nil
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return err
	}
	return nil
}
