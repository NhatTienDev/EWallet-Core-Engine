package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

type WalletHandler struct {
	walletUC domain.WalletUseCase
}

func NewWalletHandler(walletUC domain.WalletUseCase) *WalletHandler {
	return &WalletHandler{
		walletUC: walletUC,
	}
}

func (h *WalletHandler) RegisterWalletRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/api/v1/wallets", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", h.HandleCreateUserWallet)
		r.Get("/{id}", h.HandleGetUserWallet)
		r.Delete("/{id}", h.HandleDeleteUserWallet)
		r.Post("/transfer", h.HandleTransferMoney)
		r.Get("/{id}/transfers", h.HandleGetTransferHistory)
		r.Get("/{id}/entries", h.HandleGetEntryHistory)
	})
}

func getIntQueryParam(r *http.Request, key string, defaultValue int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}