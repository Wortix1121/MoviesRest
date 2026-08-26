package repository

import (
	"appMove/internal/model"
	"context"
)

type MovieRepository interface {
	Create(ctx context.Context, movie *model.Movie) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Movie, error)
	Update(ctx context.Context, movie *model.Movie) error
	Delete(ctx context.Context, id int64) error
	GetByGenre(ctx context.Context, genreID int64) ([]model.Movie, error)
	GetSeasons(ctx context.Context, movieID int64) ([]model.Season, error)
}
