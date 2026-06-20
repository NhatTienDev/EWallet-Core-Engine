package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/nhattiendev/ewallet/internal/wallet/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUserWallet(t *testing.T) {
	tests := []struct {
		name        string
		walletID    int64
		userID      int64
		setupMocks  func(repo *MockWalletRepository)
		expectedErr string
	}{
		{
			name:     "Success - delete wallet with zero balance",
			walletID: 1,
			userID:   10,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(1)).Return(&domain.Wallet{ID: 1, UserID: 10, Balance: 0, Currency: "VND"}, nil)
				repo.On("DeleteWalletByID", mock.Anything, int64(1), int64(10)).Return(nil)
			},
			expectedErr: "",
		},
		{
			name:     "Failure - wallet not found returns ErrWalletNotFound",
			walletID: 2,
			userID:   10,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(2)).Return(nil, errors.New("not found"))
			},
			expectedErr: domain.ErrWalletNotFound.Error(),
		},
		{
			name:     "Failure - forbidden access returns ErrForbiddenAccess",
			walletID: 3,
			userID:   10,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(3)).Return(&domain.Wallet{ID: 3, UserID: 11, Balance: 0, Currency: "VND"}, nil)
			},
			expectedErr: domain.ErrForbiddenAccess.Error(),
		},
		{
			name:     "Failure - remaining balance returns ErrWalletHasRemainingBalance",
			walletID: 4,
			userID:   10,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(4)).Return(&domain.Wallet{ID: 4, UserID: 10, Balance: 100, Currency: "VND"}, nil)
			},
			expectedErr: domain.ErrWalletHasRemainingBalance.Error(),
		},
		{
			name:     "Failure - delete repository error returns underlying error",
			walletID: 5,
			userID:   10,
			setupMocks: func(repo *MockWalletRepository) {
				repo.On("GetWalletByID", mock.Anything, int64(5)).Return(&domain.Wallet{ID: 5, UserID: 10, Balance: 0, Currency: "VND"}, nil)
				repo.On("DeleteWalletByID", mock.Anything, int64(5), int64(10)).Return(errors.New("delete failed"))
			},
			expectedErr: "delete failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(MockWalletRepository)
			tc.setupMocks(repo)

			uc := &walletUseCase{walletRepo: repo}
			err := uc.DeleteUserWallet(context.Background(), tc.walletID, tc.userID)
			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
			}

			repo.AssertExpectations(t)
		})
	}
}
