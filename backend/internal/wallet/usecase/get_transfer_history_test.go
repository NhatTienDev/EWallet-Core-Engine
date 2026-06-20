package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTransferHistory(t *testing.T) {
	tests := []struct {
		name        string
		walletID    int64
		limit       int32
		offset      int32
		setupMocks  func(repo *MockWalletRepository)
		expectedErr string
		expectedLen int
	}{
		{
			name:     "Success - return transfer history",
			walletID: 1,
			limit:    10,
			offset:   0,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(1)).Return(&domain.Wallet{ID: 1, UserID: 7, Balance: 1000, Currency: "VND"}, nil)
				repo.On("GetListTransfers", mock.Anything, int64(1), int32(10), int32(0)).Return([]domain.Transfer{
					{ID: 2, FromWalletID: 1, ToWalletID: 3, Amount: 500},
					{ID: 3, FromWalletID: 4, ToWalletID: 1, Amount: 200},
				}, nil)
			},
			expectedErr: "",
			expectedLen: 2,
		},
		{
			name:     "Failure - wallet not found",
			walletID: 2,
			limit:    10,
			offset:   0,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(2)).Return(nil, domain.ErrWalletNotFound)
			},
			expectedErr: domain.ErrWalletNotFound.Error(),
			expectedLen: 0,
		},
		{
			name:     "Failure - list transfer error",
			walletID: 3,
			limit:    5,
			offset:   1,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(3)).Return(&domain.Wallet{ID: 3, UserID: 9, Balance: 500, Currency: "USD"}, nil)
				repo.On("GetListTransfers", mock.Anything, int64(3), int32(5), int32(1)).Return(nil, errors.New("db error"))
			},
			expectedErr: "db error",
			expectedLen: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockWalletRepository)
			tc.setupMocks(repo)

			uc := &walletUseCase{walletRepo: repo}
			transfers, err := uc.GetTransferHistory(context.Background(), tc.walletID, tc.limit, tc.offset)

			if tc.expectedErr == "" {
				assert.NoError(t, err)
				assert.Len(t, transfers, tc.expectedLen)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				assert.Nil(t, transfers)
			}

			repo.AssertExpectations(t)
		})
	}
}
