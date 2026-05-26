package delivery

import (
	"errors"
	"net/http"

	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/internal/user/domain"
)

// @Summary     Get user profile
// @Tags        Users
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} apiResponse{data=domain.User} "User profile retrieved successfully"
// @Failure     401 {object} apiResponse "Unauthorized: Invalid user ID in context"
// @Failure     404 {object} apiResponse "User not found"
// @Failure     500 {object} apiResponse "Failed to get user profile"
// @Router      /api/v1/users/profile [get]
func (h *UserHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userIDContext := r.Context().Value(middleware.UserIDKey)

	userID, ok := userIDContext.(int64)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, apiResponse{Error: "Unauthorized: Invalid user ID in context"})
		return
	}

	user, err := h.userUC.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeJSON(w, http.StatusNotFound, apiResponse{Error: "User not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, apiResponse{Error: "Failed to get user profile"})
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Message: "User profile retrieved successfully",
		Data: user,
	})
}