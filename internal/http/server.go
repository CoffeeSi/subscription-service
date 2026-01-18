package http

import (
	"github.com/CoffeeSi/subscription-service/internal/http/handler"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
}

func NewServer(subService *service.SubscriptionService) *Server {
	router := gin.Default()

	subHandler := handler.NewSubscriptionHandler(subService)

	// Subscriptions handlers
	router.GET("/subscriptions", subHandler.ListSubscription)
	router.POST("/subscriptions", subHandler.CreateSubscription)
	router.GET("/subscriptions/:id", subHandler.GetSubscriptionsByID)
	router.DELETE("/subscriptions/:id", subHandler.DeleteSubscriptionsByID)

	router.GET("/subscriptions/totalprice", subHandler.GetTotalPrice)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{
		router: router,
	}
}

func (s *Server) Run() {
	s.router.Run()
}
