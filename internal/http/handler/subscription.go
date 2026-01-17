package handler

import (
	"net/http"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(SubscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: SubscriptionService,
	}
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var sub model.Subscription
	if err := c.BindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sub.ServiceName == "" || sub.Price <= 0 || sub.UserID == "" || sub.StartDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid fields"})
		return
	}

	id, err := h.service.CreateSubscription(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *SubscriptionHandler) ListSubscription(c *gin.Context) {
	subs, err := h.service.ListSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) GetSubscriptionsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	subs, err := h.service.GetSubscription(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) DeleteSubscriptionsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	err := h.service.DeleteSubscription(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *SubscriptionHandler) GetTotalPrice(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if userID == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, start_date and end_date parameters are required"})
		return
	}
	totalprice, err := h.service.GetTotalPriceByID(userID, serviceName, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_price": totalprice})
}
