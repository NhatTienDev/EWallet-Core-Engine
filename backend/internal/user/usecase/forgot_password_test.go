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

func TestForgotPassword(t *testing.T) {
	tests := []struct {
		name            string
		email           string
		setupMocks      func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{})
		expectedErr     error
		expectEmailSent bool
	}{
		{
			name:  "Success - create reset record and send email",
			email: "test@example.com",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				userRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&domain.User{
					ID:    10,
					Email: "test@example.com",
				}, nil)
				userRepo.On("CreatePasswordReset", mock.Anything, mock.MatchedBy(func(p *domain.PasswordReset) bool {
					return p.UserID == 10 && p.HashedToken != "" && !p.ExpiresAt.IsZero()
				})).Return(nil)
				mailRepo.On("SendResetPasswordEmail", "test@example.com", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					select {
					case emailSent <- struct{}{}:
					default:
					}
				})
			},
			expectedErr:     nil,
			expectEmailSent: true,
		},
		{
			name:  "Success - unknown email returns nil",
			email: "unknown@example.com",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				userRepo.On("GetUserByEmail", mock.Anything, "unknown@example.com").Return(nil, domain.ErrUserNotFound)
			},
			expectedErr: nil,
		},
		{
			name:  "Failure - repository error returns internal server error",
			email: "error@example.com",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				userRepo.On("GetUserByEmail", mock.Anything, "error@example.com").Return(nil, errors.New("db down"))
			},
			expectedErr: domain.ErrInternalServerError,
		},
		{
			name:  "Failure - create reset record error returns internal server error",
			email: "test@example.com",
			setupMocks: func(userRepo *MockUserRepository, mailRepo *MockMailpitSenderRepository, emailSent chan<- struct{}) {
				userRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&domain.User{
					ID:    11,
					Email: "test@example.com",
				}, nil)
				userRepo.On("CreatePasswordReset", mock.Anything, mock.AnythingOfType("*domain.PasswordReset")).Return(errors.New("insert failed"))
			},
			expectedErr: domain.ErrInternalServerError,
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

			err := uc.ForgotPassword(context.Background(), tc.email)
			assert.ErrorIs(t, err, tc.expectedErr)

			if tc.expectEmailSent {
				select {
				case <-emailSent:
					// goroutine email send called
				case <-time.After(100 * time.Millisecond):
					t.Fatal("expected SendResetPasswordEmail to be called")
				}
			}

			mockUserRepo.AssertExpectations(t)
			mockMailRepo.AssertExpectations(t)
		})
	}
}