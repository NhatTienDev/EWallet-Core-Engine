package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func respondWithJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	return respondWithJSON(w, status, map[string]string{"error": message})
}

func getIntQueryParam(r *http.Request, key string, defaultValue int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}