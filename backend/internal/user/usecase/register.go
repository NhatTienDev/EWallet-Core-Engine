package usecase

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"github.com/nhattiendev/ewallet/internal/user/domain"
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
	err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	go func(id int64) {
		u.userCreatedChan <- id
	}(user.ID)

	return user, nil
}