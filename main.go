package main

import (
	"context"
	"os"

	_ "github.com/CoffeeSi/subscription-service/docs"
	"github.com/CoffeeSi/subscription-service/internal/http"
	"github.com/CoffeeSi/subscription-service/internal/postgres"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/joho/godotenv"
)

// @title Subscription Service API
// @version 1.0
// @description API documentation for the Subscription Service.
// @host localhost:8080
// @BasePath /
func main() {
	godotenv.Load()
	ctx := context.Background()

	database := postgres.NewDatabase(ctx, os.Getenv("DB_URL"))
	defer database.Connect.Close(ctx)

	subRepo := postgres.NewSubscriptionRepository(database)
	subService := service.NewSubscriptionService(subRepo)

	server := http.NewServer(subService)
	server.Run()
}
