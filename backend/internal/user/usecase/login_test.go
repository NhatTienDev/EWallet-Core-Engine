package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nhattiendev/ewallet/internal/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		setupMocks    func(m *MockUserRepository)
		expectedErr   error
		expectToken   bool
		expectedID    int64
		expectedEmail string
	}{
		{
			name:     "Success - valid credentials returns JWT",
			email:    "success@example.com",
			password: "Password123",
			setupMocks: func(m *MockUserRepository) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.MinCost)
				m.On("GetUserByEmail", mock.Anything, "success@example.com").Return(&domain.User{
					ID:             1,
					Email:          "success@example.com",
					HashedPassword: string(hash),
				}, nil)
			},
			expectedErr:   nil,
			expectToken:   true,
			expectedID:    1,
			expectedEmail: "success@example.com",
		},
		{
			name:     "Failure - wrong password returns invalid credentials",
			email:    "success@example.com",
			password: "WrongPassword",
			setupMocks: func(m *MockUserRepository) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.MinCost)
				m.On("GetUserByEmail", mock.Anything, "success@example.com").Return(&domain.User{
					ID:             1,
					Email:          "success@example.com",
					HashedPassword: string(hash),
				}, nil)
			},
			expectedErr: domain.ErrInvalidCredentials,
			expectToken: false,
		},
		{
			name:     "Failure - email not found returns invalid credentials",
			email:    "missing@example.com",
			password: "Password123",
			setupMocks: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "missing@example.com").Return(nil, domain.ErrUserNotFound)
			},
			expectedErr: domain.ErrInvalidCredentials,
			expectToken: false,
		},
		{
			name:     "Failure - repository error returns internal server error",
			email:    "error@example.com",
			password: "Password123",
			setupMocks: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "error@example.com").Return(nil, errors.New("database down"))
			},
			expectedErr: domain.ErrInternalServerError,
			expectToken: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tc.setupMocks(mockRepo)

			uc := &userUseCase{
				userRepo:        mockRepo,
				jwtSecretKey:    []byte("test-secret"),
				userCreatedChan: nil,
			}

			token, err := uc.Login(context.Background(), tc.email, tc.password)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				claims := jwt.MapClaims{}
				parsed, parseErr := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret"), nil
				})
				assert.NoError(t, parseErr)
				assert.True(t, parsed.Valid)
				assert.Equal(t, float64(tc.expectedID), claims["user_id"])
				assert.Equal(t, tc.expectedEmail, claims["email"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}