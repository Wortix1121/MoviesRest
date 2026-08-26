package postgres

import (
	"appMove/internal/model"
	"appMove/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type watchHistoryRepository struct {
	db *sqlx.DB
}

func NewWatchHistoryRepository(db *sqlx.DB) repository.WatchHistoryRepository {
	return &watchHistoryRepository{db: db}

}

func (r *watchHistoryRepository) Upsert(ctx context.Context, entry *model.WatchHistory) error {
	const op = "Repository.Postgres.watchHistoryUpsert"
	query := `
		INSERT INTO user_watch_history (user_id, content_type, content_id, progress_seconds, watched_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, content_type, content_id) 
		DO UPDATE SET progress_seconds = EXCLUDED.progress_seconds, watched_at = NOW()       
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		entry.UserID,
		entry.ContentType,
		entry.ContentID,
		entry.ProgressSeconds,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *watchHistoryRepository) GetByUserAndContent(ctx context.Context, userID int64, contentType string, contentID int64) (*model.WatchHistory, error) {
	const op = "Repository.Postgres.GetByUserAndContent"
	query := `
		SELECT id, user_id, content_type, content_id, progress_seconds, watched_at
		FROM user_watch_history
		WHERE user_id = $1, content_type = $2, content_id = $3
	`

	var whistory model.WatchHistory
	err := r.db.GetContext(ctx, &whistory, query, userID, contentType, contentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &whistory, nil
}

func (r *watchHistoryRepository) GetRecentByUser(ctx context.Context, userID int64, limit int) ([]model.WatchHistory, error) {
	const op = "Repository.Postgres.GetRecentByUser"
	query := `
		SELECT id, user_id, content_type, content_id, progress_seconds, watched_at
		FROM user_watch_history
		WHERE user_id = $1
		ORDER BY watched_at DESC
		LIMIT $2
	`

	var whistory []model.WatchHistory
	err := r.db.SelectContext(ctx, &whistory, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return whistory, nil
}

func (r *watchHistoryRepository) Delete(ctx context.Context, id int64) error {
	const op = "Repository.Postgres.watchHistoryDelete"

	query := `
		DELETE FROM user_watch_history
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: history with id %d not found", op, id)
	}

	return nil
}
