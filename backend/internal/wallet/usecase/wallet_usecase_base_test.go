package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/stretchr/testify/mock"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWalletByID(ctx context.Context, id int64) (*domain.Wallet, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) GetWalletByIDForUpdate(ctx context.Context, id int64) (*domain.Wallet, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) UpdateWalletBalance(ctx context.Context, walletID int64, amount int64) (*domain.Wallet, error) {
	args := m.Called(ctx, walletID, amount)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Wallet), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) DeleteWalletByID(ctx context.Context, id int64, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockWalletRepository) CreateTransfer(ctx context.Context, transfer *domain.Transfer) error {
	args := m.Called(ctx, transfer)
	return args.Error(0)
}

func (m *MockWalletRepository) GetListTransfers(ctx context.Context, walletID int64, limit, offset int32) ([]domain.Transfer, error) {
	args := m.Called(ctx, walletID, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Transfer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) CreateEntry(ctx context.Context, entry *domain.Entry) error {
	args := m.Called(ctx, entry)
	return args.Error(0)
}

func (m *MockWalletRepository) GetListEntries(ctx context.Context, walletID int64, limit, offset int32) ([]domain.Entry, error) {
	args := m.Called(ctx, walletID, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Entry), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockWalletRepository) ExecTx(ctx context.Context, f func(domain.WalletRepository) error) error {
	args := m.Called(ctx, f)
	if err := args.Error(0); err != nil {
		return err
	}
	if f == nil {
		return nil
	}
	return f(m)
}
