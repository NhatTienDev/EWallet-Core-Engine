package domain

import (
	"context"
	"errors"
	"time"
)

type Wallet struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Balance int64 `json:"balance"`
	Currency string `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transfer struct {
	ID int64 `json:"id"`
	FromWalletID int64 `json:"from_wallet_id"`
	ToWalletID int64 `json:"to_wallet_id"`
	Amount int64 `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID int64 `json:"id"`
	WalletID int64 `json:"wallet_id"`
	Amount int64 `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrWalletNotFound = errors.New("Wallet not found")
	ErrInsufficientBalance = errors.New("Insufficient balance")
	ErrInvalidAmount = errors.New("Invalid amount")
	ErrSelfTransfer = errors.New("Cannot transfer to the same wallet")
	ErrWalletAlreadyExists = errors.New("Wallet already exists for this user")
	ErrInternalServerError = errors.New("Internal server error")
	ErrCurrencyMismatch = errors.New("Currency mismatch: wallets must use the same currency")
	ErrWalletHasRemainingBalance = errors.New("Cannot delete wallet with remaining balance. Please withdraw or transfer first")
	ErrForbiddenAccess = errors.New("You do not have permission to access or modify this wallet")
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *Wallet) error
	GetWalletByID(ctx context.Context, id int64) (*Wallet, error)
	GetWalletByIDForUpdate(ctx context.Context, id int64) (*Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID int64, amount int64) (*Wallet, error)
	DeleteWalletByID(ctx context.Context, id int64, userID int64) error

	CreateTransfer(ctx context.Context, transfer *Transfer) error
	GetListTransfers(ctx context.Context, walletID int64, limit, offset int32) ([]Transfer, error)

	CreateEntry(ctx context.Context, entry *Entry) error
	GetListEntries(ctx context.Context, walletID int64, limit, offset int32) ([]Entry, error)

	// The inter-table transaction execution helper function
	ExecTx(ctx context.Context, f func(WalletRepository) error) error
}

type WalletUseCase interface {
	CreateUserWallet(ctx context.Context, userID int64, currency string) (*Wallet, error)
	GetUserWallet(ctx context.Context, walletID int64) (*Wallet, error)
	DeleteUserWallet(ctx context.Context, walletID int64, userID int64) error

	// Core feature: P2P money transfer with ACID property
	TransferMoney(ctx context.Context, fromWalletID, toWalletID, amount int64) (*Transfer, error)

	GetTransferHistory(ctx context.Context, walletID int64, limit, offset int32) ([]Transfer, error)
	GetEntryHistory(ctx context.Context, walletID int64, limit, offset int32) ([]Entry, error)
}