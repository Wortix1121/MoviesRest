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

type movieRepository struct {
	db *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) repository.MovieRepository {
	return &movieRepository{db: db}

}

func (r *movieRepository) Create(ctx context.Context, movie *model.Movie) (int64, error) {
	const op = "Repository.Postgres.movieCreate"

	query := `
	INSERT INTO movies (title, description, duration_minutes, release_date, poster_url, age_rating, type, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		movie.Title,
		movie.Description,
		movie.DurationMinutes,
		movie.ReleaseDate,
		movie.PosterUrl,
		movie.AgeRating,
		movie.Type,
		movie.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: failed to create movie: %w", op, err)
	}
	return id, nil
}

func (r *movieRepository) GetByID(ctx context.Context, id int64) (*model.Movie, error) {
	const op = "Repository.Postgres.movieGetByID"

	query := `
		SELECT id, title, description, duration_minutes, release_date, poster_url, age_rating, type, created_at 
		FROM movies 
		WHERE id = $1 
	`

	var movie model.Movie
	err := r.db.GetContext(ctx, &movie, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &movie, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *model.Movie) error {
	const op = "Repository.Postgres.movieUpdate"

	query := `
		UPDATE movies
		SET title = $1, description = $2, duration_minutes = $3, release_date = $4, poster_url = $5, age_rating = $5, type = $6, created_at = $7
		WHERE id = $8
	`

	res, err := r.db.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.DurationMinutes,
		movie.ReleaseDate,
		movie.PosterUrl,
		movie.AgeRating,
		movie.Type,
		movie.CreatedAt,
		movie.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: movie with id %d not found", op, movie.ID)
	}

	return nil
}

func (r *movieRepository) Delete(ctx context.Context, id int64) error {
	const op = "Repository.Postgres.movieDelete"

	query := `
		DELETE FROM movies
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
		return fmt.Errorf("%s: movie with id %d not found", op, id)
	}

	return nil
}

func (r *movieRepository) GetByGenre(ctx context.Context, genreID int64) ([]model.Movie, error) {
	const op = "Repository.Postgres.GetByGenre"

	query := `
		SELECT m.id, m.title, m.description, m.duration_minutes, m.release_date, m.poster_url, m.age_rating, m.type, m.created_at
		FROM movies m
		JOIN movie_genres mg ON m.id = mg.movie_id
		WHERE mg.genre_id = $1
	`

	var movies []model.Movie
	err := r.db.SelectContext(ctx, &movies, query, genreID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get movies by genre: %w", op, err)
	}

	return movies, nil
}

func (r *movieRepository) GetSeasons(ctx context.Context, movieID int64) ([]model.Season, error) {
	const op = "Repository.Postgres.GetSeasons"

	query := `
		SELECT id, season_number, title, movie_id
		FROM seasons
		WHERE movie_id = $1
	`

	var seasons []model.Season
	err := r.db.SelectContext(ctx, &seasons, query, movieID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get seasons by movie id: %w", op, err)
	}
	return seasons, nil
}
