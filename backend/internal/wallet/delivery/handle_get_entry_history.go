package delivery

import (
	"net/http"
	"errors"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/respond"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func (h *WalletHandler) HandleGetEntryHistory(w http.ResponseWriter, r *http.Request) {
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

	// Only show the entries from a wallet they own
	wallet, err := h.walletUC.GetUserWallet(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			respond.WriteErrorJSON(w, http.StatusNotFound, err.Error())
			return
		}
		respond.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if wallet.UserID != authUserID {
		respond.WriteErrorJSON(w, http.StatusForbidden, "You do not have permission to view this statement")
		return
	}

	limit := getIntQueryParam(r, "limit", 10)
	offset := getIntQueryParam(r, "offset", 0)

	entries, err := h.walletUC.GetEntryHistory(r.Context(), walletID, int32(limit), int32(offset))
	if err != nil {
		respond.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	respond.WriteSuccessJSON(w, http.StatusOK, "Get entry history successfully", entries)
}