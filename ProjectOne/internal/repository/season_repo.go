package repository

import (
	"appMove/internal/model"
	"context"
)

type SeasonRepository interface {
	Create(ctx context.Context, season *model.Season) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Season, error)
	Update(ctx context.Context, season *model.Season) error
	Delete(ctx context.Context, id int64) error
}
