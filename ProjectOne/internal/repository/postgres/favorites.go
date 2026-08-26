package postgres

import (
	"appMove/internal/model"
	"appMove/internal/repository"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type favoritesRepository struct {
	db *sqlx.DB
}

func NewFavoritesRepository(db *sqlx.DB) repository.FavoritesRepository {
	return &favoritesRepository{db: db}

}

func (r *favoritesRepository) Add(ctx context.Context, userID, movieID int64) error {
	const op = "Repository.Postgres.favoritesAdd"

	query := `
	INSERT INTO user_favorites (user_id, movie_id) VALUES ($1, $2) 
	ON CONFLICT (user_id, movie_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, movieID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Удалить из избранного
func (r *favoritesRepository) Remove(ctx context.Context, userID, movieID int64) error {
	const op = "Repository.Postgres.favoritesRemove"

	query := `
	DELETE FROM user_favorites
	WHERE user_id = $1 AND movie_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, userID, movieID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Все избранные фильмы у пользователя
func (r *favoritesRepository) GetByUser(ctx context.Context, userID int64, limit, offset int) ([]model.Movie, error) {
	const op = "Repository.Postgres.favoritesGetByUser"

	query := `
	SELECT m.id, m.title, m.poster_url, m.release_date, m.duration_minutes, m.age_rating, m.type, m.created_at
	FROM movies m
	JOIN user_favorites uf ON m.id = uf.movie_id
	WHERE uf.user_id = $1
	ORDER BY uf.movbie_id
	LIMIT $2 OFFSET $3
	`

	var favorites []model.Movie
	err := r.db.SelectContext(ctx, &favorites, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return favorites, nil
}

// Проверка есть ли фильм в избранных
func (r *favoritesRepository) IsFavorite(ctx context.Context, userID, movieID int64) (bool, error) {
	const op = "Repository.Postgres.favoritesIsFavorite"

	query := `SELECT EXISTS (SELECT 1 FROM user_favorites WHERE user_id = $1 AND movie_id = $2)	`

	var isfavorites bool
	err := r.db.QueryRowContext(ctx, query, userID, movieID).Scan(&isfavorites)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isfavorites, nil
}

// Кол-во фильмов у пользователя
func (r *favoritesRepository) CountByUser(ctx context.Context, userID int64) (int, error) {
	const op = "Repository.Postgres.favoritesCountByUser"

	query := `
	SELECT COUNT(*)
	FROM user_favorites
	WHERE user_id = $1
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
