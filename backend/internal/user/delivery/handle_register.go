package delivery

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/nhattiendev/ewallet/internal/user/domain"
)

// DTO
// Data structure that client request
type registerRequest struct {
	FullName string `json:"full_name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	// Read JSON from request body
	if err := json.NewDecoder(r.Body).Decode((&req)); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "Invalid JSON format"})
		return
	}

	// Basic validation
	if req.FullName == "" || req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "Missing required fields"})
		return
	}

	user, err := h.userUC.Register(r.Context(), req.FullName, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExist) {
			writeJSON(w, http.StatusConflict, apiResponse{Error: "Email already exists"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: "Failed to register user"})
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Message: "User registered successfully",
		Data: user,
	})
}