package migrator

import (
	"appMove/pkg/config"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	srcDriver source.Driver // Драйвер источника миграций.
}

func RunMigrations(cfg *config.StorageConfig) (*Migrator, error) {
	const op = "migartor.RunMigrations"

	conStr := cfg.PathNotSSL

	migrationsPath := "file://./migrations"

	// Источник миграции
	srcDriver, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open source driver: %s, %w", op, err)
	}

	println(srcDriver)

	// Экземпляр мигратора
	m, err := migrate.NewWithSourceInstance("file", srcDriver, conStr)
	if err != nil {
		return nil, fmt.Errorf("error when migration")
	}

	defer m.Close()

	// Запускаем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("error migration up: %v", err)
	}

	log.Println("migrate completed")
	return &Migrator{srcDriver: srcDriver}, nil

}
