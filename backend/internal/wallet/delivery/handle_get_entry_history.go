package delivery

import (
	"net/http"
	"errors"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/response"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

// @Summary     Get entry history (statement) for a specific wallet
// @Tags        Wallets
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path int true "Wallet ID"
// @Param       limit query int false "Number of records to return (default 10)"
// @Param       offset query int false "Number of records to skip (default 0)"
// @Router      /api/v1/wallets/{id}/entries [get]
func (h *WalletHandler) HandleGetEntryHistory(w http.ResponseWriter, r *http.Request) {
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

	// Only show the entries from a wallet they own
	wallet, err := h.walletUC.GetUserWallet(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			response.WriteErrorJSON(w, http.StatusNotFound, err.Error())
			return
		}
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if wallet.UserID != authUserID {
		response.WriteErrorJSON(w, http.StatusForbidden, "You do not have permission to view this statement")
		return
	}

	limit := getIntQueryParam(r, "limit", 10)
	offset := getIntQueryParam(r, "offset", 0)

	entries, err := h.walletUC.GetEntryHistory(r.Context(), walletID, int32(limit), int32(offset))
	if err != nil {
		response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.WriteSuccessJSON(w, http.StatusOK, "Get entry history successfully", entries)
}