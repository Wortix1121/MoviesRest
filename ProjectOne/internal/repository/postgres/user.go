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

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}

}

func (r *userRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	const op = "Repository.Postgres.userCreate"

	query := `
	INSERT INTO users (email, password_hash, name, role, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: failed to create user: %w", op, err)
	}
	return id, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	const op = "Repository.Postgres.userGetByID"

	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at 
		FROM users 
		WHERE id = $1 
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "Repository.Postgres.userGetByEmail"

	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	const op = "Repository.Postgres.userUpdate"

	query := `
		UPDATE users
		SET  email = $1, password_hash = $2, name = $3, role = $4, updated_at = $5
		WHERE id = $6
	`

	res, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to get rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: user with id %d not found", op, user.ID)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	const op = "Repository.Postgres.userDelete"

	query := `
		DELETE FROM users
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
