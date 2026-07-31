package repository

import (
	"appMove/internal/model"
	"context"
)

type GenreRepository interface {
	Create(ctx context.Context, genre *model.Genre) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Genre, error)
	Update(ctx context.Context, genre *model.Genre) error
	Delete(ctx context.Context, id int64) error
	AddGenreToMovie(ctx context.Context, movieID int64, genreID int64) error
	RemoveGenreFromMovie(ctx context.Context, movieID int64, genreID int64) error
	GetGenresByMovie(ctx context.Context, movieID int64) ([]model.Genre, error)
	GetAll(ctx context.Context) ([]model.Genre, error)
	GetByName(ctx context.Context, name string) (*model.Genre, error)
}
