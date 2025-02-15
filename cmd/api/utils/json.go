package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ReceiveJSON[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func RespondJSON(body any, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func RespondError(err error, status int, w http.ResponseWriter) {
	errorWrapper := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}
	RespondJSON(errorWrapper, status, w)
}
