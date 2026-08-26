package repository

import (
	"appMove/internal/model"
	"context"
)

type FavoritesRepository interface {

	// Добавить в избранное
	Add(ctx context.Context, userID, movieID int64) error

	// Удалить из избранного
	Remove(ctx context.Context, userID, movieID int64) error

	// Все избранные фильмы у пользователя
	GetByUser(ctx context.Context, userID int64, limit, offset int) ([]model.Movie, error)

	//Проверка есть ли фильм в избранных
	IsFavorite(ctx context.Context, userID, movieID int64) (bool, error)

	//Кол-во фильмов у пользователя
	CountByUser(ctx context.Context, userID int64) (int, error)
}
