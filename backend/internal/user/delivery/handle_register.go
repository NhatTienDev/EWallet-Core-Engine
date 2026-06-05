package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/nhattiendev/ewallet/response"
	"github.com/nhattiendev/ewallet/internal/user/domain"
)

// DTO
// Data structure that client request
type registerRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// @Summary		Register a new user account
// @Tags 		Users
// @Accept 		json
// @Produce 	json
// @Param 		request body registerRequest true "Account registration information"
// @Success 	201 {object} response.APIResponse "User registered successfully"
// @Failure 	400 {object} response.APIResponse "Input data error (Invalid JSON format or missing required fields)"
// @Failure 	409 {object} response.APIResponse "Email already exists"
// @Failure 	500 {object} response.APIResponse "Internal server error"
// @Router 		/api/v1/users/register [post]
func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	// Read JSON from request body
	if err := json.NewDecoder(r.Body).Decode((&req)); err != nil {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	req.FullName = strings.Join(strings.Fields(req.FullName), " ")
	req.Email = strings.TrimSpace(req.Email)

	// Basic validation
	if req.FullName == "" || req.Email == "" || req.Password == "" {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	if !emailRegex.MatchString(req.Email) {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	if len(req.Password) < 10 {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Password must be at least 10 characters long")
		return
	}

	if strings.ContainsAny(req.Password, " \t\n\r") {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Password must not contain spaces")
		return
	}

	// Check at least 1 special character
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?\\"
	if !strings.ContainsAny(req.Password, specialChars) {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Password must contain at least one special character")
		return
	}

	user, err := h.userUC.Register(r.Context(), req.FullName, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.WriteErrorJSON(w, http.StatusConflict, "Email already exists")
			return
		}
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	response.WriteSuccessJSON(w, http.StatusCreated, "User registered successfully", user)
}