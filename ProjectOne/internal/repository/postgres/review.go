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

type reviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) repository.RewiewRepository {
	return &reviewRepository{db: db}

}

func (r *reviewRepository) Create(ctx context.Context, review *model.Review) (int64, error) {
	const op = "Repository.Postgres.reviewCreate"

	query := `
	INSERT INTO user_reviews (user_id, movie_id, rating, comment, created_at) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		review.UserID,
		review.MovieID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: failed to create review: %w", op, err)
	}
	return id, nil
}

func (r *reviewRepository) GetByID(ctx context.Context, id int64) (*model.Review, error) {
	const op = "Repository.Postgres.reviewGetByID"

	query := `
		SELECT id, user_id, movie_id, rating, comment, created_at 
		FROM user_reviews 
		WHERE id = $1 
	`

	var review model.Review
	err := r.db.GetContext(ctx, &review, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &review, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *model.Review) error {
	const op = "Repository.Postgres.reviewrUpdate"

	query := `
		UPDATE user_reviews
		SET  rating = $1, comment = $2
		WHERE id = $3
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		review.Rating,
		review.Comment,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: review with id %d not found", op, review.ID)
	}

	return nil
}

func (r *reviewRepository) Delete(ctx context.Context, id int64) error {
	const op = "Repository.Postgres.reviewDelete"

	query := `
		DELETE FROM user_reviews
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
		return fmt.Errorf("%s: review with id %d not found", op, id)
	}

	return nil
}

// Список отзывов с пагинацией
func (r *reviewRepository) GetReviewsByMovie(ctx context.Context, movieID int64, limit, offset int) ([]model.Review, error) {
	const op = "Repository.Postgres.reviewGetReviewsByMovie"

	query := `
		SELECT id, user_id, movie_id, rating, comment, created_at
		FROM user_reviews
		WHERE movie_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var review []model.Review
	err := r.db.SelectContext(ctx, &review, query, movieID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return review, nil
}

// Отзыв одного пользователя с проверкой
func (r *reviewRepository) GetByUserAndMovie(ctx context.Context, userID, movieID int64) (*model.Review, error) {
	const op = "Repository.Postgres.reviewGetByUserAndMovie"

	query := `
		SELECT id, user_id, movie_id, rating, comment, created_at
		FROM user_reviews
		WHERE user_id = $1 AND movie_id = $2
	`

	var review model.Review
	err := r.db.GetContext(ctx, &review, query, userID, movieID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &review, nil
}

// Средний рейтинг
func (r *reviewRepository) AverageRatingMovie(ctx context.Context, movieID int64) (float64, error) {
	const op = "Repository.Postgres.reviewAverageRatingMovie"

	query := `
		SELECT COALESCE(AVG(rating), 0)
		FROM user_reviews
		WHERE movie_id = $1
	`

	var avg float64
	err := r.db.QueryRowContext(ctx, query, movieID).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return avg, nil
}
