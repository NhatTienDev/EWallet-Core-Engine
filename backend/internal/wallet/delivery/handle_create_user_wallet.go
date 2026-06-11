package delivery

import (
	"net/http"
	"encoding/json"
	"errors"

	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/response"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

type createUserWalletRequest struct {
	Currency string `json:"currency"`
}

// @Summary     Create a new wallet for the authenticated user
// @Tags        Wallets
// @Accept      json
// @Produce     json
// @Security	BearerAuth
// @Param       request body createUserWalletRequest true "Wallet creation information"
// @Router      /api/v1/wallets [post]
func (h *WalletHandler) HandleCreateUserWallet(w http.ResponseWriter, r *http.Request) {
	// Get uerID safely from token through context
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		response.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: User ID not found in context")
		return
	}

	var req createUserWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Use authUserID to create wallet for the authenticated user
	wallet, err := h.walletUC.CreateUserWallet(r.Context(), authUserID, req.Currency)
	if err != nil {
		if errors.Is(err, domain.ErrWalletAlreadyExists) {
			response.WriteErrorJSON(w, http.StatusConflict, err.Error())
			return
		}
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.WriteSuccessJSON(w, http.StatusCreated, "Wallet created successfully", wallet)
}