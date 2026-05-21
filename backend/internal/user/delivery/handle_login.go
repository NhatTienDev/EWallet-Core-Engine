package delivery

import (
	"encoding/json"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "Invalid JSON format"})
		return
	}

	// Call to UseCase layer to get token
	token, err := h.userUC.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, apiResponse{Error: err.Error()}) // Invalid email or password
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Message: "Login successful",
		Data: map[string]string{
			"access_token": token,
		},
	})
}