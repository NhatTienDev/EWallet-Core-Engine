package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/response"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

// @Summary     Get details of a specific wallet
// @Tags        Wallets
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Wallet ID"
// @Success     200 {object} response.APIResponse{data=domain.Wallet} "Get wallet successfully"
// @Failure     400 {object} response.APIResponse "Invalid wallet ID"
// @Failure     401 {object} response.APIResponse "Unauthorized"
// @Failure     403 {object} response.APIResponse "Forbidden: You do not have permission to access this wallet"
// @Failure     404 {object} response.APIResponse "Wallet not found"
// @Failure     500 {object} response.APIResponse "Internal server error"
// @Router      /api/v1/wallets/{id} [get]
func (h *WalletHandler) HandleGetUserWallet(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: User ID not found in context")
		return
	}

	idStr := chi.URLParam(r, "id")
	walletID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	wallet, err := h.walletUC.GetUserWallet(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			response.WriteErrorJSON(w, http.StatusNotFound, err.Error())
			return
		}
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Check if this wallet belongs to the logged-in user
	if wallet.UserID != authUserID {
		response.WriteErrorJSON(w, http.StatusForbidden, "You do not have permission to access this wallet")
		return
	}

	response.WriteSuccessJSON(w, http.StatusOK, "Get wallet successfully", wallet)
}