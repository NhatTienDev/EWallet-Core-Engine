package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func (u *walletUseCase) CreateUserWallet(ctx context.Context, userID int64, currency string) (*domain.Wallet, error) {
	if currency == "" {
		currency = "VND"
	}

	// Create a new wallet with zero balance
	newWallet := &domain.Wallet{
		UserID: userID,
		Balance: 0,
		Currency: currency,
	}

	err := u.walletRepo.CreateWallet(ctx, newWallet)
	if err != nil {
		return nil, err
	}

	return newWallet, nil
}