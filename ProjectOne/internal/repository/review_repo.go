package repository

import (
	"appMove/internal/model"
	"context"
)

type RewiewRepository interface {
	Create(ctx context.Context, review *model.Review) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Review, error)
	Update(ctx context.Context, review *model.Review) error
	Delete(ctx context.Context, id int64) error

	//Список отзывов с пагинацией
	GetReviewsByMovie(ctx context.Context, movieID int64, limit, offset int) ([]model.Review, error)

	//Отзыв одного пользователя с проверкой
	GetByUserAndMovie(ctx context.Context, userID, movieID int64) (*model.Review, error)

	//Средний рейтинг
	AverageRatingMovie(ctx context.Context, movieID int64) (float64, error)
}
