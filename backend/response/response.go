package response

import (
	"net/http"
	"encoding/json"
)

type apiResponse struct {
	Error string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
}

// Write a cleaner JSON response for the client
func WriteSuccessJSON(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := apiResponse{
		Message: message,
		Data: data,
	}

	json.NewEncoder(w).Encode(res)
}

func WriteErrorJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := apiResponse{
		Message: message,
	}

	json.NewEncoder(w).Encode(res)
}