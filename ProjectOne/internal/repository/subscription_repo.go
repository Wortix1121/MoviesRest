package repository

import (
	"appMove/internal/model"
	"context"
	"time"
)

type SubscriptionRepository interface {
	// Тарифы (обычно только чтение, админка может добавлять/обновлять)
	GetPlans(ctx context.Context) ([]model.SubscriptionPlan, error)
	GetPlanByID(ctx context.Context, id int64) (*model.SubscriptionPlan, error)

	// Подписки пользователей
	Create(ctx context.Context, sub *model.Subscriptions) (int64, error)

	// активная подписка
	GetActiveByUser(ctx context.Context, userID int64) (*model.Subscriptions, error)

	// отменить (установить is_active = false)
	Cancel(ctx context.Context, userID int64) error

	// продлить (обновить end_date и is_active = true)
	Renew(ctx context.Context, userID int64, newEndDate time.Time) error
}
