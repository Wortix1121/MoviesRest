package repository

import (
	"appMove/internal/model"
	"context"
)

type EpisodeRepository interface {
	Create(ctx context.Context, episode *model.Episode) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Episode, error)
	Update(ctx context.Context, episode *model.Episode) error
	Delete(ctx context.Context, id int64) error
	GetBySeasonID(ctx context.Context, seasonID int64) ([]model.Episode, error)
	GetBySeasonAndNumber(ctx context.Context, seasonID int64, episodeNumber int) (*model.Episode, error)
	CountBySeason(ctx context.Context, seasonID int64) (int, error)
}
