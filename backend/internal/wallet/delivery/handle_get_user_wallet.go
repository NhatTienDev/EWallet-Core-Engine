package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/respond"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func (h *WalletHandler) HandleGetUserWallet(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respond.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: User ID not found in context")
		return
	}

	idStr := chi.URLParam(r, "id")
	walletID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respond.WriteErrorJSON(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	wallet, err := h.walletUC.GetUserWallet(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			respond.WriteErrorJSON(w, http.StatusNotFound, err.Error())
			return
		}
		respond.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Check if this wallet belongs to the logged-in user
	if wallet.UserID != authUserID {
		respond.WriteErrorJSON(w, http.StatusForbidden, "You do not have permission to access this wallet")
		return
	}

	respond.WriteSuccessJSON(w, http.StatusOK, "Get wallet successfully", wallet)
}