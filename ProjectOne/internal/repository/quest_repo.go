package repository

import (
	"appMove/internal/model"
	"context"
)

type QuestRepository interface {

	// CRUD для квестов (админка)
	Create(ctx context.Context, quest *model.Quest) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Quest, error)
	Update(ctx context.Context, quest *model.Quest) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]model.Quest, error)

	// Прогресс пользователя
	GetUserProgress(ctx context.Context, userID int64) ([]model.UserQuest, error)
	UpdateProgress(ctx context.Context, userID, questID int64, progress int) error
	CompleteQuest(ctx context.Context, userID, questID int64) error
}
