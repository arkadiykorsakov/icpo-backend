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
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		Error(w, http.StatusBadRequest, "Invalid request body")
		return err
	}

	if dec.More() {
		Error(w, http.StatusBadRequest, "Request body must contain a single JSON object")
		return errors.New("body contains multiple JSON values")
	}

	return nil
}
