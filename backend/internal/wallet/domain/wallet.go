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
	ErrWalletNotFound = errors.New("wallet not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount = errors.New("invalid amount")
	ErrSelfTransfer = errors.New("cannot transfer to the same wallet")
	ErrWalletAlreadyExists = errors.New("wallet already exists for this user")
	ErrInternalServerError = errors.New("internal server error")
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *Wallet) error
	GetWalletByID(ctx context.Context, id int64) (*Wallet, error)
	GetWalletByIDForUpdate(ctx context.Context, id int64) (*Wallet, error)
	UpdateWalletBalance(ctx context.Context, walletID int64, amount int64) (*Wallet, error)

	CreateTransfer(ctx context.Context, transfer *Transfer) error
	GetListTransfers(ctx context.Context, walletID int64, limit, offset int32) ([]Transfer, error)

	CreateEntry(ctx context.Context, entry *Entry) error
	GetListEntries(ctx context.Context, walletID int64, limit, offset int32) ([]Entry, error)

	// The inter-table transaction execution helper function
	ExecTx(ctx context.Context, f func())
}

type WalletUseCase interface {
	CreateUserWallet(ctx context.Context, userID int64, currency string) (*Wallet, error)
	GetUserWallet(ctx context.Context, walletID int64) (*Wallet, error)

	// Core feature: P2P money transfer with ACID property
	Transfer(ctx context.Context, fromWalletID, toWalletID, amount int64) (*Transfer, error)

	GetTransferHistory(ctx context.Context, walletID int64, limit, offset int32) ([]Transfer, error)
	GetEntryHistory(ctx context.Context, walletID int64, limit, offset int32) ([]Entry, error)
}