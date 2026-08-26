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

type seasonRepository struct {
	db *sqlx.DB
}

func NewSeasonRepository(db *sqlx.DB) repository.SeasonRepository {
	return &seasonRepository{db: db}

}

func (r *seasonRepository) Create(ctx context.Context, season *model.Season) (int64, error) {
	const op = "Repository.Postgres.seasonCreate"

	query := `
	INSERT INTO seasons (movie_id, season_number, title) 
	VALUES ($1, $2, $3) 
	RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		season.MovieID,
		season.SeasonNumber,
		season.Title,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: failed to create season: %w", op, err)
	}
	return id, nil
}

func (r *seasonRepository) GetByID(ctx context.Context, id int64) (*model.Season, error) {
	const op = "Repository.Postgres.seasonGetByID"

	query := `
		SELECT movie_id, season_number, title
		FROM seasons 
		WHERE id = $1 
	`

	var season model.Season
	err := r.db.GetContext(ctx, &season, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &season, nil
}

func (r *seasonRepository) Update(ctx context.Context, season *model.Season) error {
	const op = "Repository.Postgres.seasonUpdate"

	query := `
		UPDATE seasons
		SET movie_id = $1, season_number = $2, title = $3
		WHERE id = $4
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		season.MovieID,
		season.SeasonNumber,
		season.Title,
		season.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: season with id %d not found", op, season.ID)
	}

	return nil
}

func (r *seasonRepository) Delete(ctx context.Context, id int64) error {
	const op = "Repository.Postgres.seasonDelete"

	query := `
		DELETE FROM seasons
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
		return fmt.Errorf("%s: season with id %d not found", op, id)
	}

	return nil
}
