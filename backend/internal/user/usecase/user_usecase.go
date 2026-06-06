package usecase

import "github.com/nhattiendev/ewallet/internal/user/domain"

type userUseCase struct {
	userRepo domain.UserRepository
	jwtSecretKey []byte // Receive from main
	userCreatedChan chan<- int64 // Channel receives user_id
}

func NewUserUseCase(userRepo domain.UserRepository, jwtSecretKey string, userCreatedChan chan<- int64) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		jwtSecretKey: []byte(jwtSecretKey),
		userCreatedChan: userCreatedChan,
	}
}