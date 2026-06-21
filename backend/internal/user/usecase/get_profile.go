package usecase

import (
	"context"

	"github.com/nhattiendev/ewallet/internal/user/domain"
)

func (u *userUseCase) GetProfile(ctx context.Context, id int64) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}