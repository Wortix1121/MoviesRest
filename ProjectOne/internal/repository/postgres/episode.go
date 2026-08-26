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

type episodeRepository struct {
	db *sqlx.DB
}

func NewEpisodeRepository(db *sqlx.DB) repository.EpisodeRepository {
	return &episodeRepository{db: db}

}

func (r *episodeRepository) Create(ctx context.Context, episode *model.Episode) (int64, error) {
	const op = "Repository.Postgres.episodeCreate"

	query := `
	INSERT INTO episodes (season_id, episode_number, title, duration_minutes, video_url) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id
	`
	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		episode.SeasonID,
		episode.EpisodeNumber,
		episode.Title,
		episode.DurationMinutes,
		episode.VideoUrl,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: failed to create episode: %w", op, err)
	}
	return id, nil
}

func (r *episodeRepository) GetByID(ctx context.Context, id int64) (*model.Episode, error) {
	const op = "Repository.Postgres.episodeGetByID"

	query := `
		SELECT id, season_id, episode_number, title, duration_minutes, video_url
		FROM episodes 
		WHERE id = $1 
	`

	var episode model.Episode
	err := r.db.GetContext(ctx, &episode, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &episode, nil
}

func (r *episodeRepository) Update(ctx context.Context, episode *model.Episode) error {
	const op = "Repository.Postgres.episodeUpdate"

	query := `
		UPDATE episodes
		SET season_id = $1, episode_number = $2, title = $3 , duration_minutes = $4, video_url = $5
		WHERE id = $6
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		episode.SeasonID,
		episode.EpisodeNumber,
		episode.Title,
		episode.DurationMinutes,
		episode.VideoUrl,
		episode.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: episode with id %d not found", op, episode.ID)
	}

	return nil
}

func (r *episodeRepository) Delete(ctx context.Context, id int64) error {
	const op = "Repository.Postgres.episodeDelete"

	query := `
		DELETE FROM episodes
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
		return fmt.Errorf("%s: episode with id %d not found", op, id)
	}

	return nil
}

func (r *episodeRepository) GetBySeasonID(ctx context.Context, seasonID int64) ([]model.Episode, error) {
	const op = "Repository.Postgres.GetBySeasonID"

	query := `
		SELECT id, seeason_id, episode_number, title, duration_minutes, video_url
		FROM episodes
		WHERE season_id = $1
		ORDER BY episode_number
	`

	var episodes []model.Episode
	err := r.db.SelectContext(ctx, &episodes, query, seasonID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return episodes, nil
}

func (r *episodeRepository) GetBySeasonAndNumber(ctx context.Context, seasonID int64, episodeNumber int) (*model.Episode, error) {
	const op = "Repository.Postgres.GetBySeasonAndNumber"

	query := `
		SELECT id, season_id, episode_number, title, duration_minutes, video_url
		FROM episodes
		WHERE season_id = $1 AND episode_number = $2
	`
	var episode model.Episode
	err := r.db.GetContext(ctx, &episode, query, seasonID, episodeNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &episode, nil
}

func (r *episodeRepository) CountBySeason(ctx context.Context, seasonID int64) (int, error) {
	const op = "Repository.Postgres.CountBySeason"

	query := `
		SELECT COUNT(*)
		FROM episodes
		WHERE season_id = $1
	`

	var count int64
	err := r.db.GetContext(ctx, &count, query, seasonID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(count), nil
}
