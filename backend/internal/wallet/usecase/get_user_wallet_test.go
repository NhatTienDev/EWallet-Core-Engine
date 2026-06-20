package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserWallet(t *testing.T) {
	tests := []struct {
		name        string
		walletID    int64
		setupMocks  func(repo *MockWalletRepository)
		expectedErr string
	}{
		{
			name:     "Success - wallet returned",
			walletID: 42,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(42)).Return(&domain.Wallet{
					ID:       42,
					UserID:   7,
					Balance:  1000,
					Currency: "VND",
				}, nil)
			},
			expectedErr: "",
		},
		{
			name:     "Failure - wallet not found",
			walletID: 99,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(99)).Return(nil, domain.ErrWalletNotFound)
			},
			expectedErr: "Wallet not found",
		},
		{
			name:     "Failure - repository error",
			walletID: 100,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(100)).Return(nil, errors.New("db error"))
			},
			expectedErr: "db error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockWalletRepository)
			tc.setupMocks(repo)

			uc := &walletUseCase{
				walletRepo: repo,
			}

			wallet, err := uc.GetUserWallet(context.Background(), tc.walletID)
			if tc.expectedErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, wallet)
				assert.Equal(t, tc.walletID, wallet.ID)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				assert.Nil(t, wallet)
			}

			repo.AssertExpectations(t)
		})
	}
}
