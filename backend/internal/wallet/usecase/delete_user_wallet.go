package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func (u *walletUseCase) DeleteUserWallet(ctx context.Context, walletID int64, userID int64) error {
	wallet, err := u.walletRepo.GetWalletByID(ctx, walletID)
	if err != nil {
		return domain.ErrWalletNotFound
	}

	if wallet.Balance > 0 {
		return domain.ErrWalletHasRemainingBalance
	}

	return u.walletRepo.DeleteWalletByID(ctx, walletID, userID)
}