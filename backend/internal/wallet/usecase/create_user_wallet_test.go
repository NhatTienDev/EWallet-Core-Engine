package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserWallet(t *testing.T) {
	tests := []struct {
		name         string
		userID       int64
		currency     string
		setupMocks   func(repo *MockWalletRepository)
		expectedErr  string
		expectedCurr string
	}{
		{
			name:     "Success - create wallet with default currency",
			userID:   10,
			currency: "",
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("CreateWallet", mock.Anything, mock.MatchedBy(func(w *domain.Wallet) bool {
					return w.UserID == 10 && w.Balance == 0 && w.Currency == "VND"
				})).Return(nil)
			},
			expectedErr:  "",
			expectedCurr: "VND",
		},
		{
			name:     "Success - create wallet with specified currency",
			userID:   11,
			currency: "USD",
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("CreateWallet", mock.Anything, mock.MatchedBy(func(w *domain.Wallet) bool {
					return w.UserID == 11 && w.Balance == 0 && w.Currency == "USD"
				})).Return(nil)
			},
			expectedErr:  "",
			expectedCurr: "USD",
		},
		{
			name:     "Failure - create wallet repository error",
			userID:   12,
			currency: "EUR",
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("CreateWallet", mock.Anything, mock.AnythingOfType("*domain.Wallet")).Return(errors.New("insert failed"))
			},
			expectedErr:  "insert failed",
			expectedCurr: "EUR",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockWalletRepository)
			tc.setupMocks(repo)

			uc := &walletUseCase{
				walletRepo: repo,
			}

			wallet, err := uc.CreateUserWallet(context.Background(), tc.userID, tc.currency)
			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			}

			if tc.expectedErr == "" {
				assert.NotNil(t, wallet)
				assert.Equal(t, tc.userID, wallet.UserID)
				assert.Equal(t, tc.expectedCurr, wallet.Currency)
				assert.Equal(t, int64(0), wallet.Balance)
			} else {
				assert.Nil(t, wallet)
			}

			repo.AssertExpectations(t)
		})
	}
}
