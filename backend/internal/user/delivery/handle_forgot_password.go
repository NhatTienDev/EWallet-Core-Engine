package delivery

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nhattiendev/ewallet/response"
)

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *UserHandler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Email is required")
		return
	}

	err := h.userUC.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.WriteSuccessJSON(w, http.StatusOK, "If the account exists, a password reset link has been sent to your email.", nil)
}