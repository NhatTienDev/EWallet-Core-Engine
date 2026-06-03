package delivery

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/nhattiendev/ewallet/middleware"
)

type transferMoneyRequest struct {
	FromWalletID int64 `json:"from_wallet_id"`
	ToWalletID int64 `json:"to_wallet_id"`
	Amount int64 `json:"amount"`
}

func (h *WalletHandler) HandleTranssferMoney(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: User ID not found in context")
		return
	}

	var req transferMoneyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if sender wallet belongs to logged-in user
	fromWallet, err := h.walletUC.GetUserWallet(r.Context(), req.FromWalletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			respondWithError(w, http.StatusNotFound, "Source wallet not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if fromWallet.UserID != authUserID {
		respondWithError(w, http.StatusForbidden, "You can only transfer money from your own wallet")
	}

	// Transfer money after authorized successfully
	transfer, err := h.walletUC.TransferMoney(r.Context(), req.FromWalletID, req.ToWalletID, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSelfTransfer), errors.Is(err, domain.ErrInvalidAmount):
			respondWithError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrInsufficientBalance):
			respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, domain.ErrWalletNotFound):
			respondWithError(w, http.StatusNotFound, "Destination wallet not found")
		default:
			respondWithError(w, http.StatusInternalServerError, "Internal server error") // "Transaction failed. Roll back."
		}
		return
	}

	respondWithJSON(w, http.StatusOK, transfer)
}