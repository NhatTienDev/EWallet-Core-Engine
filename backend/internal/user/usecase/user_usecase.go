package usecase

import "github.com/nhattiendev/ewallet/internal/user/domain"

type userUseCase struct {
	userRepo domain.UserRepository
	jwtSecretKey []byte // Receive from main
}

func NewUserUseCase(userRepo domain.UserRepository, jwtSecretKey string) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		jwtSecretKey: []byte(jwtSecretKey),
	}
}