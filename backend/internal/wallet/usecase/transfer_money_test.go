package usecase

import (
	"context"
	"testing"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransferMoney(t *testing.T) {
	tests := []struct {
		name        string
		fromID      int64
		toID        int64
		amount      int64
		setupMocks  func(repo *MockWalletRepository)
		expectedErr string
	}{
		{
			name:   "Success - transfer between same currency wallets",
			fromID: 1,
			toID:   2,
			amount: 500,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("ExecTx", mock.Anything, mock.Anything).Return(nil)
				repo.On("GetWalletByIDForUpdate", mock.Anything, int64(1)).Return(&domain.Wallet{ID: 1, UserID: 10, Balance: 1000, Currency: "VND"}, nil)
				repo.On("GetWalletByIDForUpdate", mock.Anything, int64(2)).Return(&domain.Wallet{ID: 2, UserID: 11, Balance: 100, Currency: "VND"}, nil)
				repo.On("UpdateWalletBalance", mock.Anything, int64(1), int64(-500)).Return(&domain.Wallet{ID: 1, UserID: 10, Balance: 500, Currency: "VND"}, nil)
				repo.On("CreateEntry", mock.Anything, mock.MatchedBy(func(e *domain.Entry) bool {
					return e.WalletID == 1 && e.Amount == -500
				})).Return(nil)
				repo.On("UpdateWalletBalance", mock.Anything, int64(2), int64(500)).Return(&domain.Wallet{ID: 2, UserID: 11, Balance: 600, Currency: "VND"}, nil)
				repo.On("CreateEntry", mock.Anything, mock.MatchedBy(func(e *domain.Entry) bool {
					return e.WalletID == 2 && e.Amount == 500
				})).Return(nil)
				repo.On("CreateTransfer", mock.Anything, mock.MatchedBy(func(t *domain.Transfer) bool {
					return t.FromWalletID == 1 && t.ToWalletID == 2 && t.Amount == 500
				})).Return(nil)
			},
			expectedErr: "",
		},
		{
			name:        "Failure - self transfer returns ErrSelfTransfer",
			fromID:      1,
			toID:        1,
			amount:      100,
			setupMocks:  func(repo *MockWalletRepository) {},
			expectedErr: domain.ErrSelfTransfer.Error(),
		},
		{
			name:        "Failure - invalid amount returns ErrInvalidAmount",
			fromID:      1,
			toID:        2,
			amount:      0,
			setupMocks:  func(repo *MockWalletRepository) {},
			expectedErr: domain.ErrInvalidAmount.Error(),
		},
		{
			name:   "Failure - insufficient balance returns ErrInsufficientBalance",
			fromID: 1,
			toID:   2,
			amount: 1500,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("ExecTx", mock.Anything, mock.Anything).Return(nil)
				repo.On("GetWalletByIDForUpdate", mock.Anything, int64(1)).Return(&domain.Wallet{ID: 1, UserID: 10, Balance: 1000, Currency: "VND"}, nil)
				repo.On("GetWalletByIDForUpdate", mock.Anything, int64(2)).Return(&domain.Wallet{ID: 2, UserID: 11, Balance: 100, Currency: "VND"}, nil)
			},
			expectedErr: domain.ErrInsufficientBalance.Error(),
		},
		{
			name:   "Failure - currency mismatch returns ErrCurrencyMismatch",
			fromID: 1,
			toID:   2,
			amount: 100,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("ExecTx", mock.Anything, mock.Anything).Return(nil)
				repo.On("GetWalletByIDForUpdate", mock.Anything, int64(1)).Return(&domain.Wallet{ID: 1, UserID: 10, Balance: 1000, Currency: "VND"}, nil)
				repo.On("GetWalletByIDForUpdate", mock.Anything, int64(2)).Return(&domain.Wallet{ID: 2, UserID: 11, Balance: 100, Currency: "USD"}, nil)
			},
			expectedErr: domain.ErrCurrencyMismatch.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockWalletRepository)
			tc.setupMocks(repo)

			uc := &walletUseCase{walletRepo: repo}
			transfer, err := uc.TransferMoney(context.Background(), tc.fromID, tc.toID, tc.amount)

			if tc.expectedErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, transfer)
				assert.Equal(t, tc.fromID, transfer.FromWalletID)
				assert.Equal(t, tc.toID, transfer.ToWalletID)
				assert.Equal(t, tc.amount, transfer.Amount)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				assert.Nil(t, transfer)
			}

			repo.AssertExpectations(t)
		})
	}
}
