package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/nhattiendev/ewallet/internal/user/domain"
)

func (u *userUseCase) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Hash the plaintext token received from the client to match it with the hash stored in DB
	hash := sha256.Sum256([]byte(token))
	hashedTokenString := fmt.Sprintf("%x", hash)

	// Get reset token record from DB
	resetRecord, err := u.userRepo.GetValidPasswordReset(ctx, hashedTokenString)
	if err != nil {
		return err // Wrong token, used token, expired token
	}

	user, err := u.userRepo.GetUserByID(ctx, resetRecord.UserID)
	if err != nil {
		return domain.ErrInternalServerError
	}

	// Encode new password with bcrypt
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return domain.ErrInternalServerError
	}

	// Update new user password to DB
	err = u.userRepo.UpdateUserPassword(ctx, user.ID, string(newHashedPassword))
	if err != nil {
		return domain.ErrInternalServerError
	}

	// Token is for single use only
	err = u.userRepo.MarkPasswordResetUsed(ctx, resetRecord.ID)
	if err != nil {
		return domain.ErrInternalServerError
	}

	// Send password change alert email through Background Job (Goroutine)
	go func() {
		err := u.mailpitSenderRepo.SendPasswordChangedAlert(user.Email)
		if err != nil {
			fmt.Printf("[Worker] Error sending password change alert email for %s: %v", user.Email, err)
		}
	}()

	return nil
}