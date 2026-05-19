package usecase

import "github.com/nhattiendev/ewallet/internal/user/domain"

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCse(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}