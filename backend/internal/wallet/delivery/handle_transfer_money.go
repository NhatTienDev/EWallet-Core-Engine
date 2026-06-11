package delivery

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/response"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

type transferMoneyRequest struct {
	FromWalletID int64 `json:"from_wallet_id"`
	ToWalletID int64 `json:"to_wallet_id"`
	Amount int64 `json:"amount"`
}

// @Summary     Transfer money between wallets
// @Tags        Wallets
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body transferMoneyRequest true "Transfer details (FromWalletID, ToWalletID, Amount)"
// @Router      /api/v1/wallets/transfer [post]
func (h *WalletHandler) HandleTransferMoney(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: User ID not found in context")
		return
	}

	var req transferMoneyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if sender wallet belongs to logged-in user
	fromWallet, err := h.walletUC.GetUserWallet(r.Context(), req.FromWalletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			response.WriteErrorJSON(w, http.StatusNotFound, "Wallet not found")
			return
		}
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if fromWallet.UserID != authUserID {
		response.WriteErrorJSON(w, http.StatusForbidden, "You can only transfer money from your own wallet")
		return
	}

	// Transfer money after authorized successfully
	transfer, err := h.walletUC.TransferMoney(r.Context(), req.FromWalletID, req.ToWalletID, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSelfTransfer), errors.Is(err, domain.ErrInvalidAmount), errors.Is(err, domain.ErrCurrencyMismatch):
			response.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrInsufficientBalance):
			response.WriteErrorJSON(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, domain.ErrWalletNotFound):
			response.WriteErrorJSON(w, http.StatusNotFound, "Destination wallet not found")
		default:
			response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error") // "Transaction failed. Roll back."
		}
		return
	}

	response.WriteSuccessJSON(w, http.StatusOK, "Money transferred successfully", transfer)
}