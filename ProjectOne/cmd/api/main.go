package main

import (
	"appMove/internal/storage/postgre"
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
	defer store.Close()

	// Проверка здоровья БД
	if err := store.HealthCheck(ctx); err != nil {
		log.Fatal("Database health check failed:", err)
	}

	fmt.Println("✅ Successfully connected to database!")

	//init http - (gin)

	//init middleware

	//init server

	//start server

}
