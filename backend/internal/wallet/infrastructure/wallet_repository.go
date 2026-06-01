package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

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

func mapToEntryDomain(dbEntry sqlc.Entry) domain.Entry {
	return domain.Entry{
		ID: dbEntry.ID,
		WalletID: dbEntry.WalletID,
		Amount: dbEntry.Amount,
		CreatedAt: dbEntry.CreatedAt,
	}
}

// ExecTx func implements a block of code (callback) inside DB transaction
func (r *walletRepository) ExecTx(ctx context.Context, f func(domain.WalletRepository) error) error {
    // Open a new transaction
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    // WithTx func allows queries to run on transaction instead of connection pool
    qTx := r.q.WithTx(tx)
    txRepo := &walletRepository{
        db: r.db,
        q: qTx,
    }

	// Check result and Commit/Rollback transaction
    err = f(txRepo)
    if err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
        return err
    }

    return tx.Commit()
}

func (r *walletRepository) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	arg := sqlc.CreateWalletParams{
		UserID: wallet.UserID,
		Balance: wallet.Balance,
		Currency: wallet.Currency,
	}

	result, err := r.q.CreateWallet(ctx, arg)
	if err != nil {
		return err
	}

	*wallet = mapToWalletDomain(result)
	return nil
}

func (r *walletRepository) GetWalletByID(ctx context.Context, id int64) (*domain.Wallet, error) {
	result, err := r.q.GetWalletByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}

	wallet := mapToWalletDomain(result)
	return &wallet, nil
}

func (r *walletRepository) GetWalletByIDForUpdate(ctx context.Context, id int64) (*domain.Wallet, error) {
	result, err := r.q.GetWalletByIDForUpdate(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrWalletNotFound
		}
		return nil, err
	}

	wallet := mapToWalletDomain(result)
	return &wallet, nil
}

func (r *walletRepository) UpdateWalletBalance(ctx context.Context, walletID int64, amount int64) (*domain.Wallet, error) {
	arg := sqlc.AddWalletBalanceParams{	
		ID: walletID,
		Balance: amount,
	}

	result, err := r.q.AddWalletBalance(ctx, arg)
	if err != nil {
		return nil, err
	}

	wallet := mapToWalletDomain(result)
	return &wallet, nil
}

func (r *walletRepository) CreateTransfer(ctx context.Context, transfer *domain.Transfer) error {
	arg := sqlc.CreateTransferParams{
		FromWalletID: transfer.FromWalletID,
		ToWalletID: transfer.ToWalletID,
		Amount: transfer.Amount,
	}

	result, err := r.q.CreateTransfer(ctx, arg)
	if err != nil {
		return err
	}

	*transfer = mapToTransferDomain(result)
	return nil
}

func (r *walletRepository) GetListTransfers(ctx context.Context, walletID int64, limit, offset int32) ([]domain.Transfer, error) {
	arg := sqlc.GetListTransfersParams{
		FromWalletID: walletID,
		Limit: limit,
		Offset: offset,
	}

	result, err := r.q.GetListTransfers(ctx, arg)
	if err != nil {
		return nil, err
	}

	transfers := make([]domain.Transfer, len(result))
	for i, t := range result {
		transfers[i] = mapToTransferDomain(t)
	}

	return transfers, nil
}

func (r *walletRepository) CreateEntry(ctx context.Context, entry *domain.Entry) error {
	arg := sqlc.CreateEntryParams{
		WalletID: entry.WalletID,
		Amount: entry.Amount,
	}

	result, err := r.q.CreateEntry(ctx, arg)
	if err != nil {
		return err
	}

	*entry = mapToEntryDomain(result)
	return nil
}

func (r *walletRepository) GetListEntries(ctx context.Context, walletID int64, limit, offset int32) ([]domain.Entry, error) {
	arg := sqlc.GetListEntriesParams{
		WalletID: walletID,
		Limit: limit,
		Offset: offset,
	}

	result, err := r.q.GetListEntries(ctx, arg)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.Entry, len(result))
	for i, e := range result {
		entries[i] = mapToEntryDomain(e)
	}

	return entries, nil
}