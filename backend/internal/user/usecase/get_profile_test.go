package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/nhattiendev/ewallet/internal/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name        string
		id          int64
		setupMocks  func(m *MockUserRepository)
		expectedErr error
	}{
		{
			name: "Success - user profile returned",
			id:   1,
			setupMocks: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, int64(1)).Return(&domain.User{
					ID:       1,
					FullName: "Test User",
					Email:    "test@example.com",
				}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "Failure - user not found returns repository error",
			id:   2,
			setupMocks: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, int64(2)).Return(nil, domain.ErrUserNotFound)
			},
			expectedErr: domain.ErrUserNotFound,
		},
		{
			name: "Failure - repository error returns internal error",
			id:   3,
			setupMocks: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, int64(3)).Return(nil, errors.New("db unavailable"))
			},
			expectedErr: errors.New("db unavailable"),
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

			user, err := uc.GetProfile(context.Background(), tc.id)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.id, user.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}