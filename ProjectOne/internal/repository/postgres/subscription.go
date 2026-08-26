package postgres

import (
	"appMove/internal/model"
	"appMove/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type subscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) repository.SubscriptionRepository {
	return &subscriptionRepository{db: db}

}

// Тарифы (обычно только чтение, админка может добавлять/обновлять)
func (r *subscriptionRepository) GetPlans(ctx context.Context) ([]model.SubscriptionPlan, error) {
	const op = "Repository.Postgres.SubsriptionGetPlans"

	query := `
		SELECT id, name, price, features
		FROM subscription_plans
	`

	var plans []model.SubscriptionPlan
	err := r.db.SelectContext(ctx, &plans, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return plans, nil
}

func (r *subscriptionRepository) GetPlanByID(ctx context.Context, id int64) (*model.SubscriptionPlan, error) {
	const op = "Repository.Postgres.SubsriptionGetPlanByID"

	query := `
		SELECT id, name, price, features
		FROM subscription_plans
		WHERE id = $1
	`

	var plan model.SubscriptionPlan
	err := r.db.GetContext(ctx, &plan, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &plan, nil
}

// Подписки пользователей
func (r *subscriptionRepository) Create(ctx context.Context, sub *model.Subscriptions) (int64, error) {
	const op = "Repository.Postgres.SubsriptionCreate"

	query := `
		INSERT INTO subscriptions (user_id, plan_id, start_date, end_date, is_active)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int64
	err := r.db.QueryRowContext(ctx, query,
		sub.UserID,
		sub.PlanID,
		sub.StartDate,
		sub.EndDate,
		sub.IsActive,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create subsription: %w", op, err)
	}

	return id, nil
}

// активная подписка
func (r *subscriptionRepository) GetActiveByUser(ctx context.Context, userID int64) (*model.Subscriptions, error) {
	const op = "Repository.Postgres.SubsriptionGetActiveByUser"

	query := `
		SELECT id, user_id, plan_id, start_date, end_date, is_active
		FROM subscriptions
		WHERE user_id = $1 AND is_active = true
	`
	var sub model.Subscriptions
	err := r.db.GetContext(ctx, &sub, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &sub, nil
}

// отменить (установить is_active = false)
func (r *subscriptionRepository) Cancel(ctx context.Context, userID int64) error {
	const op = "Repository.Postgres.SubsriptionCancel"

	query := `
		UPDATE subscriptions SET is_active = false
		WHERE user_id = $1 AND is_active = true
	`

	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowAffected == 0 {
		return fmt.Errorf("%s: active subscription for user %d not found", op, userID)

	}

	return nil
}

// продлить (обновить end_date и is_active = true)
func (r *subscriptionRepository) Renew(ctx context.Context, userID int64, newEndDate time.Time) error {
	const op = "Repository.Postgres.SubsriptionRenew"

	query := `
		UPDATE subscriptions SET end_date = $1, is_active = true
		WHERE user_id = $2 AND is_active = true
	`

	res, err := r.db.ExecContext(ctx, query, newEndDate, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowAffected == 0 {
		return fmt.Errorf("%s: active subscription for user %d not found", op, userID)

	}

	return nil
}
