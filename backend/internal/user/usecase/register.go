package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
)

func (u *userUseCase) Register(ctx context.Context, fullName, email, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	user := &domain.User{
		FullName: fullName,
		Email:    email,
		HashedPassword: string(hashedPassword),
	}

	// Call Infrastructure layer to save to PostgreSQL
	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}