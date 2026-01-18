package http

import (
	"github.com/CoffeeSi/subscription-service/internal/http/handler"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	router *gin.Engine
	Logger *zap.Logger
}

func NewServer(subService *service.SubscriptionService, logger *zap.Logger) *Server {
	router := gin.Default()

	subHandler := handler.NewSubscriptionHandler(subService, logger)

	// Subscriptions handlers
	router.GET("/subscriptions", subHandler.ListSubscription)
	router.POST("/subscriptions", subHandler.CreateSubscription)
	router.GET("/subscriptions/:id", subHandler.GetSubscriptionsByID)
	router.DELETE("/subscriptions/:id", subHandler.DeleteSubscriptionsByID)

	router.GET("/subscriptions/totalprice", subHandler.GetTotalPrice)

	return &Server{
		router: router,
		Logger: logger,
	}
}

func (s *Server) Run() {
	s.router.Run()
}
