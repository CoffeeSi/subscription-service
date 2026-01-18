package handler

import (
	"net/http"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
	logger  *zap.Logger
}

func NewSubscriptionHandler(SubscriptionService *service.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: SubscriptionService,
		logger:  logger,
	}
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var sub model.Subscription
	if err := c.BindJSON(&sub); err != nil {
		h.logger.Error("invalid subscription creation request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sub.ServiceName == "" || sub.Price <= 0 || sub.UserID == "" || sub.StartDate == "" {
		h.logger.Warn("missing or invalid fields in subscription creation request",
			zap.String("subscription", sub.UserID),
			zap.String("service_name", sub.ServiceName),
			zap.Int64("price", sub.Price),
			zap.String("start_date", sub.StartDate),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid fields"})
		return
	}

	id, err := h.service.CreateSubscription(sub)
	if err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("subscription created", zap.Int("id", id))
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *SubscriptionHandler) ListSubscription(c *gin.Context) {
	subs, err := h.service.ListSubscriptions()
	if err != nil {
		h.logger.Error("failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) GetSubscriptionsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Warn("id parameter is required for getting subscription", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	subs, err := h.service.GetSubscription(id)
	if err != nil {
		h.logger.Warn("failed to get subscriptions by id", zap.String("user_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) DeleteSubscriptionsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Warn("id parameter is required for deletion", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	err := h.service.DeleteSubscription(id)
	if err != nil {
		h.logger.Warn("failed to delete subscriptions", zap.String("user_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("subscriptions deleted", zap.String("user_id", id))
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *SubscriptionHandler) GetTotalPrice(c *gin.Context) {
	h.logger.Info("get subscriptions total price for a period",
		zap.String("user_id", c.Query("user_id")),
		zap.String("service_name", c.Query("service_name")),
		zap.String("start_date", c.Query("start_date")),
		zap.String("end_date", c.Query("end_date")),
	)
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if userID == "" || startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, start_date and end_date parameters are required"})
		return
	}
	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	totalprice, err := h.service.GetTotalPriceByID(userID, serviceNamePtr, startDate, endDate)
	if err != nil {
		h.logger.Error("failed to get total price",
			zap.String("user_id", userID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("total price calculated", zap.Int64("total_price", totalprice))
	c.JSON(http.StatusOK, gin.H{"total_price": totalprice})
}
