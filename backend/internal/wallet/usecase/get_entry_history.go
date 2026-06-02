package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func (u *walletUseCase) GetEntryHistory(ctx context.Context, walletID int64, limit, offset int32) ([]domain.Entry, error) {
	_, err := u.walletRepo.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return u.walletRepo.GetListEntries(ctx, walletID, limit, offset)
}