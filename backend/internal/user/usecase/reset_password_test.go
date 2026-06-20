package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nhattiendev/ewallet/internal/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResetPassword(t *testing.T) {
	tests := []struct {
		name            string
		token           string
		newPassword     string
		setupMocks      func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{})
		expectedErr     error
		expectAlertSent bool
	}{
		{
			name:        "Success - valid token updates password and sends alert",
			token:       "valid-token",
			newPassword: "NewPassword123",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				hashedToken := "397a2a9c5bf5e2ccec38c2596b682bb1bd05fe6e4ecea6c10cf42755ff225403" // sha256("valid-token")
				userRepo.On("GetValidPasswordReset", mock.Anything, hashedToken).Return(&domain.PasswordReset{
					ID:     100,
					UserID: 10,
				}, nil)
				userRepo.On("GetUserByID", mock.Anything, int64(10)).Return(&domain.User{
					ID:    10,
					Email: "test@example.com",
				}, nil)
				userRepo.On("UpdateUserPassword", mock.Anything, int64(10), mock.AnythingOfType("string")).Return(nil)
				userRepo.On("MarkPasswordResetUsed", mock.Anything, int64(100)).Return(nil)
				mailRepo.On("SendPasswordChangedAlert", "test@example.com").Return(nil).Run(func(args mock.Arguments) {
					select {
					case emailSent <- struct{}{}:
					default:
					}
				})
			},
			expectedErr:     nil,
			expectAlertSent: true,
		},
		{
			name:        "Failure - invalid reset token returns error",
			token:       "invalid-token",
			newPassword: "NewPassword123",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				hashedToken := "644d0c3b82bfe5e0665a116b2eb139d6abd6c9083ede0891237b1723e0010a14" // sha256("invalid-token")
				userRepo.On("GetValidPasswordReset", mock.Anything, hashedToken).Return(nil, domain.ErrInvalidResetToken)
			},
			expectedErr:     domain.ErrInvalidResetToken,
			expectAlertSent: false,
		},
		{
			name:        "Failure - user lookup error returns internal server error",
			token:       "valid-token",
			newPassword: "NewPassword123",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				hashedToken := "397a2a9c5bf5e2ccec38c2596b682bb1bd05fe6e4ecea6c10cf42755ff225403"
				userRepo.On("GetValidPasswordReset", mock.Anything, hashedToken).Return(&domain.PasswordReset{
					ID:     100,
					UserID: 10,
				}, nil)
				userRepo.On("GetUserByID", mock.Anything, int64(10)).Return(nil, errors.New("db down"))
			},
			expectedErr:     domain.ErrInternalServerError,
			expectAlertSent: false,
		},
		{
			name:        "Failure - update password error returns internal server error",
			token:       "valid-token",
			newPassword: "NewPassword123",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				hashedToken := "397a2a9c5bf5e2ccec38c2596b682bb1bd05fe6e4ecea6c10cf42755ff225403"
				userRepo.On("GetValidPasswordReset", mock.Anything, hashedToken).Return(&domain.PasswordReset{
					ID:     100,
					UserID: 10,
				}, nil)
				userRepo.On("GetUserByID", mock.Anything, int64(10)).Return(&domain.User{
					ID:    10,
					Email: "test@example.com",
				}, nil)
				userRepo.On("UpdateUserPassword", mock.Anything, int64(10), mock.AnythingOfType("string")).Return(errors.New("update fail"))
			},
			expectedErr:     domain.ErrInternalServerError,
			expectAlertSent: false,
		},
		{
			name:        "Failure - mark reset used error returns internal server error",
			token:       "valid-token",
			newPassword: "NewPassword123",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				hashedToken := "397a2a9c5bf5e2ccec38c2596b682bb1bd05fe6e4ecea6c10cf42755ff225403"
				userRepo.On("GetValidPasswordReset", mock.Anything, hashedToken).Return(&domain.PasswordReset{
					ID:     100,
					UserID: 10,
				}, nil)
				userRepo.On("GetUserByID", mock.Anything, int64(10)).Return(&domain.User{
					ID:    10,
					Email: "test@example.com",
				}, nil)
				userRepo.On("UpdateUserPassword", mock.Anything, int64(10), mock.AnythingOfType("string")).Return(nil)
				userRepo.On("MarkPasswordResetUsed", mock.Anything, int64(100)).Return(errors.New("mark failed"))
			},
			expectedErr:     domain.ErrInternalServerError,
			expectAlertSent: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(MockUserRepository)
			mockMailRepo := new(MockMailpitSenderRepository)
			emailSent := make(chan struct{}, 1)
			tc.setupMocks(mockUserRepo, mockMailRepo, emailSent)

			uc := &userUseCase{
				userRepo:          mockUserRepo,
				mailpitSenderRepo: mockMailRepo,
				jwtSecretKey:      []byte("test-secret"),
				userCreatedChan:   nil,
			}

			err := uc.ResetPassword(context.Background(), tc.token, tc.newPassword)
			assert.ErrorIs(t, err, tc.expectedErr)

			if tc.expectAlertSent {
				select {
				case <-emailSent:
					// alert sent
				case <-time.After(100 * time.Millisecond):
					t.Fatal("expected SendPasswordChangedAlert to be called")
				}
			}

			mockUserRepo.AssertExpectations(t)
			mockMailRepo.AssertExpectations(t)
		})
	}
}
