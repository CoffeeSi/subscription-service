package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/CoffeeSi/subscription-service/docs"
	"github.com/CoffeeSi/subscription-service/internal/http"
	"github.com/CoffeeSi/subscription-service/internal/logging"
	"github.com/CoffeeSi/subscription-service/internal/postgres"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title Subscription Service API
// @version 1.0
// @description API documentation for the Subscription Service.
// @host localhost:8080
// @BasePath /
func main() {
	godotenv.Load()

	logger, _ := logging.NewLogging()

	ctx := context.Background()
	DB_URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	database := postgres.NewDatabase(ctx, DB_URL)
	defer database.Connect.Close(ctx)

	m, err := migrate.New(
		"file://migrations",
		DB_URL,
	)
	if err != nil {
		logger.Fatal("Failed to start migration", zap.Error(err))
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("Migration failed", zap.Error(err))
	}

	subRepo := postgres.NewSubscriptionRepository(database, logger)
	subService := service.NewSubscriptionService(subRepo, logger)

	server := http.NewServer(subService, logger)
	server.Run()
}
