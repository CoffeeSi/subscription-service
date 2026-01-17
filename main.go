package main

import (
	"context"
	"os"

	"github.com/CoffeeSi/subscription-service/internal/http"
	"github.com/CoffeeSi/subscription-service/internal/postgres"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	database := postgres.NewDatabase(ctx, os.Getenv("DB_URL"))
	defer database.Connect.Close(ctx)

	subRepo := postgres.NewSubscriptionRepository(database)
	subHandler := service.NewSubscription(subRepo)

	server := http.NewServer(*subHandler)
	server.Run()
}
