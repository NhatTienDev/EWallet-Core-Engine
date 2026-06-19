package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/nhattiendev/ewallet/response"
	"github.com/nhattiendev/ewallet/internal/user/domain"
)

type resetPasswordRequest struct {
	Token string `json:"token"`
	NewPassword string `json:"new_password"`
}

// @Summary      Reset to a new password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        request body resetPasswordRequest true "Reset token and new password"
// @Router       /api/v1/users/reset-password [post]
func (h *UserHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	req.Token = strings.TrimSpace(req.Token)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	
	if req.Token == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Reset token is required")
		return
	}

	if len(req.NewPassword) < 10 {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Password must be at least 10 characters long")
		return
	}

	err := h.userUC.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidResetToken) {
			response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid or expired reset token")
			return
		}

		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.WriteSuccessJSON(w, http.StatusOK, "Your password has been reset successfully. Please login with your new password.", nil)
}