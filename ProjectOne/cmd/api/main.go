package main

import (
	"appMove/internal/storage/postgre"
	"appMove/migrator"
	"appMove/pkg/config"
	"context"
	"fmt"
	"log"
)

func main() {

	// Init config (viper)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	fmt.Println(cfg)

	// Init database (postgre)
	ctx := context.Background()

	store, err := postgre.New(ctx, &cfg.Storage)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer store.CloseDB()

	// Проверка работоспособности БД
	if err := store.HealthCheckDB(ctx); err != nil {
		log.Fatal("Database health check failed:", err)
	}

	fmt.Println("✅ Successfully connected to database!")

	// Migrator
	mRuning, err := migrator.RunMigrations(&cfg.Storage)
	if err != nil {
		log.Fatal("Failed to migrate ", err)
	}

	_ = mRuning

	//init http - (gin)

	//init middleware

	//init server

	//start server

}
