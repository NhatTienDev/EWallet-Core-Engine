package delivery

import (
	"net/http"

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
	r.Route("api/v1/wallets", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Post("/", h.HandleCreateUserWallet)
			r.Get("/{id}", h.HandleGetUserWallet)
			r.Post("/transfer", h.HandleTransferMoney)
			r.Get("/{id}/transfers", h.HandleGetTransferHistory)
			r.Get("/{id}/entries", h.HandleGetEntryHistory)
		})
	})
}