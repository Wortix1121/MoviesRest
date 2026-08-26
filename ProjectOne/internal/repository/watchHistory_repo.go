package repository

import (
	"appMove/internal/model"
	"context"
)

type WatchHistoryRepository interface {
	//Создать или обновить прогресс
	Upsert(ctx context.Context, entry *model.WatchHistory) error

	//Получение прогресса конкретного пользователя по конкретному контенту
	GetByUserAndContent(ctx context.Context, userID int64, contentType string, contentID int64) (*model.WatchHistory, error)

	//Последнии записи для пользователя
	GetRecentByUser(ctx context.Context, userID int64, limit int) ([]model.WatchHistory, error)
	Delete(ctx context.Context, id int64) error
}
