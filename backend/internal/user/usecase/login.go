package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nhattiendev/ewallet/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
)

func (u *userUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", domain.ErrInternalServerError
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Initialize payload for JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"exp": time.Now().Add(24 * time.Hour).Unix(), // expire after 24h
	}

	// Generate JWT token and  sign with HS256 algorithm + secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.jwtSecretKey)
	if err != nil {
		return "", domain.ErrInternalServerError
	}

	return tokenString, nil
}