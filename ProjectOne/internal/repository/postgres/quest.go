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

type questRepository struct {
	db *sqlx.DB
}

func NewQuestRepository(db *sqlx.DB) repository.QuestRepository {
	return &questRepository{db: db}

}

func (r *questRepository) Create(ctx context.Context, quest *model.Quest) (int64, error) {
	const op = "Repository.Postgres.questCreate"

	query := `
	INSERT INTO quests (title, description, type, required_count, reward_points) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		quest.Title,
		quest.Description,
		quest.Type,
		quest.RequiredCount,
		quest.RewardPoints,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: failed to create user: %w", op, err)
	}
	return id, nil
}

func (r *questRepository) GetByID(ctx context.Context, id int64) (*model.Quest, error) {
	const op = "Repository.Postgres.questGetByID"

	query := `
		SELECT id, title, description, type, required_count, reward_points 
		FROM quests 
		WHERE id = $1 
	`

	var quest model.Quest
	err := r.db.GetContext(ctx, &quest, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &quest, nil
}

func (r *questRepository) Update(ctx context.Context, quest *model.Quest) error {
	const op = "Repository.Postgres.questUpdate"

	query := `
		UPDATE quests
		SET  title = $1, description = $2, type = $3, required_count = $4, reward_points = $5
		WHERE id = $6
	`

	res, err := r.db.ExecContext(ctx, query,
		quest.Title,
		quest.Description,
		quest.Type,
		quest.RequiredCount,
		quest.RewardPoints,
		quest.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: quest with id %d not found", op, quest.ID)
	}

	return nil
}

func (r *questRepository) Delete(ctx context.Context, id int64) error {
	const op = "Repository.Postgres.questDelete"

	query := `
		DELETE FROM quests
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
		return fmt.Errorf("%s: quests with id %d not found", op, id)
	}

	return nil
}

func (r *questRepository) GetAll(ctx context.Context) ([]model.Quest, error) {
	const op = "Repository.Postgres.questGetAll"

	query := `
		SELECT id, title, description, type, required_count, reward_points
		FROM quests
		ORDER BY id
	`

	var quests []model.Quest
	err := r.db.SelectContext(ctx, &quests, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return quests, nil
}

func (r *questRepository) GetUserProgress(ctx context.Context, userID int64) ([]model.UserQuest, error) {
	const op = "Repository.Postgres.questGetUserProgress"

	query := `
		SELECT user_id, quest_id, progress, complete_at
		FROM user_quests
		WHERE user_id = $1
	`

	var progress []model.UserQuest
	err := r.db.SelectContext(ctx, &progress, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return progress, nil
}

func (r *questRepository) UpdateProgress(ctx context.Context, userID, questID int64, progress int) error {
	const op = "Repository.Postgres.questUpdateProgress"

	query := `
		INSERT INTO user_quests (user_id, quest_id, progress)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, quest_id)
		DO UPDATE SET progress = EXCLUDED.progress
	`

	_, err := r.db.ExecContext(ctx, query, userID, questID, progress)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *questRepository) CompleteQuest(ctx context.Context, userID, questID int64) error {
	const op = "Repository.Postgres.questCompleteQuest"

	query := `
		UPDATE user_quests
		SET complete_at = NOW()
		WHERE user_id = $1 AND quest_id = $2 AND complete_at IS NULL
	`

	res, err := r.db.ExecContext(ctx, query, userID, questID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: progress not found or already completed (user=%d, quest=%d)", op, userID, questID)
	}

	return nil
}
