package usecase

import "github.com/nhattiendev/ewallet/internal/wallet/domain"

type walletUseCase struct {
	walletRepo domain.WalletRepository
}

func NewWalletUseCase(walletRepo domain.WalletRepository) domain.WalletUseCase {
	return &walletUseCase{
		walletRepo: walletRepo,
	}
}