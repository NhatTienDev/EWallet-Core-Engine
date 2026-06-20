package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/nhattiendev/ewallet/internal/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name          string
		fullName      string
		email         string
		password      string
		setupMocks    func(m *MockUserRepository)
		expectedErr   error
		expectChannel bool
	}{
		{
			name:     "Success - New user registered successfully",
			fullName: "Nhat Tien Dev",
			email:    "nhattiendev@gmail.com",
			password: "MySecurePassword123",
			setupMocks: func(m *MockUserRepository) {
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte("MySecurePassword123"))
					return u.FullName == "Nhat Tien Dev" && u.Email == "nhattiendev@gmail.com" && err == nil
				})).Run(func(args mock.Arguments) {
					user := args.Get(1).(*domain.User)
					user.ID = 99 
				}).Return(nil)
			},
			expectedErr:   nil,
			expectChannel: true,
		},
		{
			name:     "Failure - Email already exists in system",
			fullName: "Duplicate User",
			email:    "exist@gmail.com",
			password: "Password123",
			setupMocks: func(m *MockUserRepository) {
				m.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(domain.ErrEmailAlreadyExists)
			},
			expectedErr:   domain.ErrEmailAlreadyExists,
			expectChannel: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Khởi tạo Mock Repo
			mockRepo := new(MockUserRepository)
			tc.setupMocks(mockRepo)

			userChan := make(chan int64, 1)
			uc := &userUseCase{
				userRepo:        mockRepo,
				userCreatedChan: userChan,
			}

			resultUser, err := uc.Register(context.Background(), tc.fullName, tc.email, tc.password)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, resultUser)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resultUser)
				assert.Equal(t, int64(99), resultUser.ID)
			}

			if tc.expectChannel {
				select {
				case idFromChan := <-userChan:
					assert.Equal(t, int64(99), idFromChan)
				case <-time.After(50 * time.Millisecond):
					t.Fatal("Timeout: Expected user ID was not sent to userCreatedChan")
				}
			} else {
				assert.Len(t, userChan, 0, "Channel should be empty on failure cases")
			}

			mockRepo.AssertExpectations(t)
		})
	}
}