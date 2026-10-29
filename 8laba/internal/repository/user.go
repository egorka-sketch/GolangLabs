package repository

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID       int
	Name     string
	Birthday time.Time
}

var ErrNotFound = errors.New("пользователь не найден")

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int) (User, error)
	Create(ctx context.Context, user User) (User, error)
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id int) error
}
