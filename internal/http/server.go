package http

import (
	"net/http"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router              *gin.Engine
	subscriptionHandler service.Subscription
}

func NewServer(subscriptionHandler service.Subscription) *Server {
	router := gin.Default()
	router.GET("/subscriptions", func(c *gin.Context) {
		subs, err := subscriptionHandler.ListSubscriptions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, subs)
	})

	router.POST("/subscriptions", func(c *gin.Context) {
		var sub model.Subscription
		if err := c.BindJSON(&sub); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := subscriptionHandler.CreateSubscription(sub)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	})

	router.GET("/subscriptions/:id", func(c *gin.Context) {
		id := c.Param("id")
		subs, err := subscriptionHandler.GetSubscription(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, subs)
	})

	router.DELETE("/subscriptions/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := subscriptionHandler.DeleteSubscription(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	return &Server{
		router:              router,
		subscriptionHandler: subscriptionHandler,
	}
}

func (s *Server) Run() {
	s.router.Run()
}
