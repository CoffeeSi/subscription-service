package handler

import (
	"net/http"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/CoffeeSi/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// @Summary Create a new subscription
// @Description Create a new subscription with the provided details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Subscription to create"
// @Success 201 {object} map[string]int "ID of the created subscription"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var sub model.Subscription
	if err := c.BindJSON(&sub); err != nil {
		h.logger.Error("invalid subscription creation request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sub.ServiceName == "" || sub.Price < 0 || sub.UserID == uuid.Nil || sub.StartDate == "" {
		h.logger.Warn("missing or invalid fields in subscription creation request",
			zap.String("subscription", sub.UserID.String()),
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

// @Summary List all subscriptions
// @Description Retrieve a list of all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription "List of subscriptions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscription(c *gin.Context) {
	subs, err := h.service.ListSubscriptions()
	if err != nil {
		h.logger.Error("failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// @Summary Get subscriptions by user ID
// @Description Retrieve subscriptions for a specific user by their ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} model.Subscription "List of subscriptions for the user"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscriptionsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Warn("id parameter is required for getting subscription", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("invalid UUID format for id parameter", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	subs, err := h.service.GetSubscription(userID)
	if err != nil {
		h.logger.Warn("failed to get subscriptions by id", zap.String("user_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// @Summary Update subscriptions by user ID
// @Description Update subscriptions for a specific user by their ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param service_name path string true "Service Name"
// @Param subscription body model.Subscription true "Subscription data to update"
// @Success 200 {object} map[string]string "Success update message"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions/{id}/{service_name} [put]
func (h *SubscriptionHandler) UpdateSubscriptionByID(c *gin.Context) {
	id := c.Param("id")
	serviceName := c.Param("service_name")
	if id == "" {
		h.logger.Warn("id parameter is required for getting subscription", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}
	if serviceName == "" {
		h.logger.Warn("service_name parameter is required for updating subscription", zap.String("service_name", serviceName))
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_name parameter is required"})
		return
	}

	var sub model.Subscription
	if err := c.BindJSON(&sub); err != nil {
		h.logger.Error("invalid subscription creation request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if sub.ServiceName == "" || sub.Price < 0 || sub.UserID == uuid.Nil || sub.StartDate == "" {
		h.logger.Warn("missing or invalid fields in subscription creation request",
			zap.String("subscription", sub.UserID.String()),
			zap.String("service_name", sub.ServiceName),
			zap.Int64("price", sub.Price),
			zap.String("start_date", sub.StartDate),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid fields"})
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("invalid UUID format for id parameter", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	err = h.service.UpdateSubscriptionByID(userID, serviceName, sub)
	if err != nil {
		h.logger.Warn("failed to update subscriptions", zap.String("user_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("subscriptions updated", zap.String("user_id", id))
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// @Summary Delete subscriptions by user ID
// @Description Delete subscriptions for a specific user by their ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "Success delete message"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscriptionsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Warn("id parameter is required for deleting subscription", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("invalid UUID format for id parameter", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	err = h.service.DeleteSubscription(userID)
	if err != nil {
		h.logger.Warn("failed to delete subscriptions", zap.String("user_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("subscriptions deleted", zap.String("user_id", id))
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// @Summary Get total price of subscriptions
// @Description Calculate the total price of subscriptions for a user within a date range, optionally filtered by service name
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID"
// @Param service_name query string false "Service Name"
// @Param start_date query string true "Start Date (MM-YYYY)"
// @Param end_date query string true "End Date (MM-YYYY)"
// @Success 200 {object} map[string]int64 "Total price of subscriptions"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions/totalprice [get]
func (h *SubscriptionHandler) GetTotalPrice(c *gin.Context) {
	h.logger.Info("get subscriptions total price for a period",
		zap.String("user_id", c.Query("user_id")),
		zap.String("service_name", c.Query("service_name")),
		zap.String("start_date", c.Query("start_date")),
		zap.String("end_date", c.Query("end_date")),
	)
	id := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if id == "" || startDate == "" || endDate == "" {
		h.logger.Warn("missing required parameters for total price calculation",
			zap.String("user_id", id),
			zap.String("start_date", startDate),
			zap.String("end_date", endDate),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, start_date and end_date parameters are required"})
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Warn("invalid UUID format for id parameter", zap.String("user_id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	totalprice, err := h.service.GetTotalPriceByID(userID, serviceNamePtr, startDate, endDate)
	if err != nil {
		h.logger.Error("failed to get total price",
			zap.String("user_id", id),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("total price calculated", zap.Int64("total_price", totalprice))
	c.JSON(http.StatusOK, gin.H{"total_price": totalprice})
}
