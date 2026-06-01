package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func (u *walletUseCase) GetUserWallet(ctx context.Context, walletID int64) (*domain.Wallet, error) {
	return u.walletRepo.GetWalletByID(ctx, walletID)
}