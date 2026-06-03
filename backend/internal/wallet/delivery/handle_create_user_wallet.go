package delivery

import (
	"net/http"
	"encoding/json"
	"errors"

	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

type createWalletRequest struct {
	Currency string `json:"currency"`
}

func (h *WalletHandler) HandleCreateUserWallet(w http.ResponseWriter, r *http.Request) {
	// Get uerID safely from token through context
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: User ID not found in context")
	}

	var req createWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Use authUserID to create wallet for the authenticated user
	wallet, err := h.walletUC.CreateUserWallet(r.Context(), authUserID, req.Currency)
	if err != nil {
		if errors.Is(err, domain.ErrWalletAlreadyExists) {
			respondWithError(w, http.StatusConflict, err.Error())
		}
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}

	respondWithJSON(w, http.StatusCreated, wallet)
}