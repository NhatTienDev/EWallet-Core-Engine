package delivery

import (
	"errors"
	"net/http"

	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/response"
	"github.com/nhattiendev/ewallet/internal/user/domain"
)

// @Summary     Get user profile
// @Tags        Users
// @Produce     json
// @Security	BearerAuth
// @Router      /api/v1/users/profile [get]
func (h *UserHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userIDContext := r.Context().Value(middleware.UserIDKey)

	userID, ok := userIDContext.(int64)
	if !ok {
		response.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: Invalid user ID in context")
		return
	}

	user, err := h.userUC.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.WriteErrorJSON(w, http.StatusNotFound, "User not found")
			return
		}
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	response.WriteSuccessJSON(w, http.StatusOK, "User profile retrieved successfully", user)
}