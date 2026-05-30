package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/nhattiendev/ewallet/internal/wallet/infrastructure/sqlc"
)

type walletRepository struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewWalletRepository(db *sql.DB) domain.WalletRepository {
	return &walletRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

// Convert SQLC data type to domain type
func mapToWalletDomain(dbWallet sqlc.Wallet) domain.Wallet {
	return domain.Wallet{
		ID: dbWallet.ID,
		UserID: dbWallet.UserID,
		Balance: dbWallet.Balance,
		Currency: dbWallet.Currency,
		CreatedAt: dbWallet.CreatedAt,
		UpdatedAt: dbWallet.UpdatedAt,
	}
}

func mapToTransferDomain(dbTransfer sqlc.Transfer) domain.Transfer {
	return domain.Transfer{
		ID: dbTransfer.ID,
		FromWalletID: dbTransfer.FromWalletID,
		ToWalletID: dbTransfer.ToWalletID,
		Amount: dbTransfer.Amount,
		CreatedAt: dbTransfer.CreatedAt,
	}
}