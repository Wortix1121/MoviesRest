package postgre

import (
	"appMove/pkg/config"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// Новое подключение
func New(ctx context.Context, cfg *config.StorageConfig) (*Storage, error) {
	const op = "Storage.New"

	conStr := fmt.Sprintf("%s", cfg.Path)
	db, err := sqlx.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("Failed Connected to open connection: %s, %w", op, err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxLifetime * time.Hour)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("Failed to connect to DB: %s, %w", op, err)
	}

	return &Storage{db: db}, nil

}

// Возвращаем экземпляр DB
func (s *Storage) GetDB() *sqlx.DB {
	return s.db
}

// Проверка доступности DB
func (s *Storage) HealthCheck(ctx context.Context) error {
	const op = "Storage.HealthCheck"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("DB health check failde: %s, %w", op, err)
	}

	return nil
}

// Закрываем подключение DB
func (s *Storage) Close() error {
	const op = "Storage.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("Failed to close DB connection: %s, %w", op, err)
	}
	return nil
}
