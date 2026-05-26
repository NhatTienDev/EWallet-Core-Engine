package delivery

import (
	"encoding/json"
	"net/http"
	"strings"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary     Login to the system
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       request body loginRequest true "Login information"
// @Success     200 {object} apiResponse{data=map[string]string} "Login successful"
// @Failure     400 {object} apiResponse "Invalid JSON format"
// @Failure     401 {object} apiResponse "Invalid email or password"
// @Router      /api/v1/users/login [post]
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "Invalid JSON format"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "Email and password are required"})
		return
	}

	if !emailRegex.MatchString(req.Email) {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "Invalid email or password"})
		return
	}

	if len(req.Password) > 72 {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "Invalid email or password"})
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