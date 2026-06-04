package delivery

import (
	"net/http"
	"encoding/json"
	"errors"

	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/respond"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

type createUserWalletRequest struct {
	Currency string `json:"currency"`
}

func (h *WalletHandler) HandleCreateUserWallet(w http.ResponseWriter, r *http.Request) {
	// Get uerID safely from token through context
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respond.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: User ID not found in context")
		return
	}

	var req createUserWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Use authUserID to create wallet for the authenticated user
	wallet, err := h.walletUC.CreateUserWallet(r.Context(), authUserID, req.Currency)
	if err != nil {
		if errors.Is(err, domain.ErrWalletAlreadyExists) {
			respond.WriteErrorJSON(w, http.StatusConflict, err.Error())
		}
		respond.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	respond.WriteSuccessJSON(w, http.StatusCreated, "Wallet created successfully", wallet)
}