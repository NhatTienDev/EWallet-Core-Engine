package delivery

import (
	"net/http"
	"encoding/json"
)

// Struct formats returned result
type apiResponse struct {
	Error string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
}

// Write a cleaner JSON response for the client
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}