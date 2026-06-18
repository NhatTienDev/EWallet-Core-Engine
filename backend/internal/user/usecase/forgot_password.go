package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/nhattiendev/ewallet/internal/user/domain"
)

func (u *userUseCase) ForgotPassword(ctx context.Context, email string) error {
	// Check user existence
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil // Even if the email address doesn't exist, so that hackers can't access it
		}
		return domain.ErrInternalServerError
	}

	// Generate random reset token with high entropy
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return domain.ErrInternalServerError
	}

	// Encode to Base64 URL-safe string to include in the link sent to the user
	tokenString := base64.RawURLEncoding.EncodeToString(tokenBytes)

	// Hash toke string with SHA-256 before saving to DB
	hash := sha256.Sum256([]byte(tokenString))
	hashedTokenString := fmt.Sprintf("%x", hash)

	// Set a reasonable expiration date
	expiresAt := time.Now().Add(15 * time.Minute)

	// Save reset token info to DB
	resetRecord := &domain.PasswordReset{
		UserID: user.ID,
		HashedToken: hashedTokenString,
		ExpiresAt: expiresAt,
	}

	if err := u.userRepo.CreatePasswordReset(ctx, resetRecord); err != nil {
		return domain.ErrInternalServerError
	}

	// Send email containing plaintext token through Background Job (Goroutine)
	go func() {
		err := u.mailpitSenderRepo.SendResetPasswordEmail(user.Email, tokenString)
		if err != nil {
			// In real systems, using logger (zap, logrus)
			fmt.Printf("[Worker] Error sending password reset email for %s: %v", user.Email, err)
		}
	}()

	return nil
}