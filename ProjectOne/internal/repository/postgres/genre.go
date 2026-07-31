package postgres

import (
	"appMove/internal/model"
	"appMove/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type genreRepository struct {
	db *sqlx.DB
}

func NewGenreRepository(db *sqlx.DB) repository.GenreRepository {
	return &genreRepository{db: db}

}

func (r *genreRepository) Create(ctx context.Context, genre *model.Genre) (int64, error) {
	const op = "Repository.Postgres.Create"

	query := `
	INSERT INTO genres (name) 
	VALUES ($1) 
	RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		genre.Name,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: failed to create user: %w", op, err)
	}
	return id, nil
}

func (r *genreRepository) GetByID(ctx context.Context, id int64) (*model.Genre, error) {
	const op = "Repository.Postgres.GetByID"

	query := `
		SELECT name
		FROM genres 
		WHERE id = $1 
	`

	var genre model.Genre
	err := r.db.GetContext(ctx, &genre, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &genre, nil
}

func (r *genreRepository) Update(ctx context.Context, genre *model.Genre) error {
	const op = "Repository.Postgres.Update"

	query := `
		UPDATE genres
		SET name = $1
		WHERE id = $2
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		genre.Name,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: user with id %d not found", op, genre.ID)
	}

	return nil
}

func (r *genreRepository) Delete(ctx context.Context, id int64) error {
	const op = "Repository.Postgres.Delete"

	query := `
		DELETE FROM genres
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
		return fmt.Errorf("%s: user with id %d not found", op, id)
	}

	return nil
}

func (r *genreRepository) AddGenreToMovie(ctx context.Context, movieID int64, genreID int64) error {
	const op = "Repository.Postgres.AddGenreToMovie"

	query := `
	INSERT INTO movie_genres (movie_id, genre_id) 
	VALUES ($1, $2) 
	`

	_, err := r.db.ExecContext(ctx, query, movieID, genreID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *genreRepository) RemoveGenreFromMovie(ctx context.Context, movieID int64, genreID int64) error {
	const op = "Repository.Postgres.RemoveGenreFromMovie"

	query := `
	DELETE FROM movie_genres
	WHERE movie_id = $1 AND genre_id = $2
	`

	res, err := r.db.ExecContext(ctx, query, movieID, genreID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: genre %d not attached to movie %d", op, genreID, movieID)
	}

	return nil
}

func (r *genreRepository) GetGenresByMovie(ctx context.Context, movieID int64) ([]model.Genre, error) {
	const op = "Repository.Postgres.GetGenresByMovie"

	query := `
		SELECT g.id, g.name
		FROM genres g
		JOIN movie_genres mg ON g.id = mg.genre_id
		WHERE mg.movie_id = $1 
	`

	var genre []model.Genre
	err := r.db.SelectContext(ctx, &genre, query, movieID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genre, nil
}

func (r *genreRepository) GetAll(ctx context.Context) ([]model.Genre, error) {
	const op = "Repository.Postgres.GetAll"

	query := `
		SELECT id, name
		FROM genres 
	`

	var genre []model.Genre
	err := r.db.SelectContext(ctx, &genre, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return genre, nil
}

func (r *genreRepository) GetByName(ctx context.Context, name string) (*model.Genre, error) {
	const op = "Repository.Postgres.GetByName"

	query := `
		SELECT id, name
		FROM genres 
		WHERE name = $1
	`

	var genre model.Genre
	err := r.db.GetContext(ctx, &genre, query, name)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &genre, nil
}
