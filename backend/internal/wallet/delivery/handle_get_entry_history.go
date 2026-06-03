package delivery

import (
	"net/http"
	"errors"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func (h *WalletHandler) HandleGetEntryHistory(w http.ResponseWriter, r *http.Request) {
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: User ID not found in context")
		return
	}

	idStr := chi.URLParam(r, "id")
	walletID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	// Only show the entries from a wallet they own
	wallet, err := h.walletUC.GetUserWallet(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if wallet.UserID != authUserID {
		respondWithError(w, http.StatusForbidden, "You do not have permission to view this statement")
		return
	}

	limit := getIntQueryParam(r, "limit", 10)
	offset := getIntQueryParam(r, "offset", 0)

	entries, err := h.walletUC.GetEntryHistory(r.Context(), walletID, int32(limit), int32(offset))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, entries)
}