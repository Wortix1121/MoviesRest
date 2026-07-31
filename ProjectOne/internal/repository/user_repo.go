package repository

import (
	"appMove/internal/model"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user *model.User) error
}
