package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/nhattiendev/ewallet/response"
	"github.com/nhattiendev/ewallet/internal/user/domain"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary		Login to the system
// @Tags        Users
// @Accept      json
// @Produce     json
// @Param       request body loginRequest true "Login information"
// @Success     200 {object} response.APIResponse{data=map[string]string} "Login successfully"
// @Failure     400 {object} response.APIResponse "Invalid JSON format or missing required fields"
// @Failure     401 {object} response.APIResponse "Invalid email or password"
// @Failure     500 {object} response.APIResponse "Internal server error"
// @Router      /api/v1/users/login [post]
func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	if !emailRegex.MatchString(req.Email) {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	if len(req.Password) > 72 {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	// Call to UseCase layer to get token
	token, err := h.userUC.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrUserNotFound) {
			response.WriteErrorJSON(w, http.StatusUnauthorized, "Invalid email or password") // Invalid email or password
			return
		}

		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.WriteSuccessJSON(w, http.StatusOK, "Login successfully", map[string]string{
		"access_token": token,
	})
}