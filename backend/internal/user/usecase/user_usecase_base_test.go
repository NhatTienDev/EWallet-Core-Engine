package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/user/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type MockMailpitSenderRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreatePasswordReset(ctx context.Context, p *domain.PasswordReset) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockUserRepository) GetValidPasswordReset(ctx context.Context, t string) (*domain.PasswordReset, error) {
	args := m.Called(ctx, t)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.PasswordReset), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) MarkPasswordResetUsed(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, id int64, p string) error {
	args := m.Called(ctx, id, p)
	return args.Error(0)
}

func (m *MockMailpitSenderRepository) SendResetPasswordEmail(toEmail string, token string) error {
	args := m.Called(toEmail, token)
	return args.Error(0)
}

func (m *MockMailpitSenderRepository) SendPasswordChangedAlert(toEmail string) error {
	args := m.Called(toEmail)
	return args.Error(0)
}