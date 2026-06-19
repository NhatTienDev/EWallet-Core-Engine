package domain

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

type PasswordReset struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	HashedToken string    `json:"-"`
	IsUsed      bool      `json:"is_used"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	ErrUserNotFound        = errors.New("User not found")
	ErrEmailAlreadyExists  = errors.New("Email already exists")
	ErrInternalServerError = errors.New("Internal server error")
	ErrInvalidCredentials  = errors.New("Invalid email or password")
	ErrInvalidResetToken   = errors.New("Invalid or expired reset token")
)

// Infrastructure interface
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)

	CreatePasswordReset(ctx context.Context, passwordReset *PasswordReset) error
	GetValidPasswordReset(ctx context.Context, hashedToken string) (*PasswordReset, error)
	MarkPasswordResetUsed(ctx context.Context, id int64) error
	UpdateUserPassword(ctx context.Context, userID int64, newHashedPassword string) error
}

// UseCase interface
type UserUseCase interface {
	Register(ctx context.Context, fullName, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetProfile(ctx context.Context, id int64) (*User, error)

	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}