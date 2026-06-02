package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
)

func (u *walletUseCase) TransferMoney(ctx context.Context, fromWalletID, toWalletID, amount int64) (*domain.Transfer, error) {
	if fromWalletID == toWalletID {
		return nil, domain.ErrSelfTransfer
	}

	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	var transfer domain.Transfer

	// Execute the entire sequence of operations within a DB transaction
	err := u.walletRepo.ExecTx(ctx, func(txRepo domain.WalletRepository) error {
		// Lock account that has smaller ID first to prevent deadlock
		var fromWallet *domain.Wallet
		// var toWallet *domain.Wallet
		var errGet error

		if fromWalletID < toWalletID {
			// Lock sender wallet first, then receiver wallet
			fromWallet, errGet = txRepo.GetWalletByIDForUpdate(ctx, fromWalletID)
			if errGet != nil {
				return errGet
			}

			_, errGet = txRepo.GetWalletByIDForUpdate(ctx, toWalletID)
			if errGet != nil {
				return errGet
			}
		} else {
			// Lock receiver wallet first, then sender wallet
			_, errGet = txRepo.GetWalletByIDForUpdate(ctx, toWalletID)
			if errGet != nil {
				return errGet
			}

			fromWallet, errGet = txRepo.GetWalletByIDForUpdate(ctx, fromWalletID)
			if errGet != nil {
				return errGet
			}
		}

		// Check account balance is sufficient to send money
		if fromWallet.Balance < amount {
			return domain.ErrInsufficientBalance
		}

		// Subtract money from sender wallet and make negative entry audit
		_, errGet = txRepo.UpdateWalletBalance(ctx, fromWalletID, -amount)
		if errGet != nil {
			return errGet
		}

		errGet = txRepo.CreateEntry(ctx, &domain.Entry{
			WalletID: fromWalletID,
			Amount: -amount,
		})
		if errGet != nil {
			return errGet
		}

		// Add money to receiver wallet and make positive entry audit
		_, errGet = txRepo.UpdateWalletBalance(ctx, toWalletID, amount)
		if errGet != nil {
			return errGet
		}

		errGet = txRepo.CreateEntry(ctx, &domain.Entry{
			WalletID: toWalletID,
			Amount: amount,
		})
		if errGet != nil {
			return errGet
		}

		// Create transaction proof record
		transfer = domain.Transfer{
			FromWalletID: fromWalletID,
			ToWalletID: toWalletID,
			Amount: amount,
		}

		errGet = txRepo.CreateTransfer(ctx, &transfer)
		if errGet != nil {
			return errGet
		}
		return nil // If not error, ExecTx func automatically commit the transaction
	})

	if err != nil {
		return nil, err // If error, ExecTx func automatically rollback the transaction
	}

	return &transfer, nil
}