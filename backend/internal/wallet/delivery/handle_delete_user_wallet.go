package delivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/nhattiendev/ewallet/middleware"
	"github.com/nhattiendev/ewallet/response"
)

func (h *WalletHandler) HandleDeleteUserWallet(w http.ResponseWriter, r *http.Request) {
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

	err = h.walletUC.DeleteUserWallet(r.Context(), walletID, authUserID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWalletNotFound):
			response.WriteErrorJSON(w, http.StatusNotFound, "Wallet not found")
		case errors.Is(err, domain.ErrWalletHasRemainingBalance):
			response.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		default:
			response.WriteErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	response.WriteSuccessJSON(w, http.StatusOK, "Wallet deleted successfully", nil)
}