package domain

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID int64 `json:"id"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	HashedPassword string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailAlreadyExist = errors.New("email already exists")
	ErrInternalServerError = errors.New("internal server error")
)

// Infrastructure interface
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

// Usecase interface
type UserUsecase interface {
	Register(ctx context.Context, fullName, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, error)
	GetProfile(ctx context.Context, id int64) (*User, error)
}