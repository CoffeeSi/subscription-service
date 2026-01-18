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

// @Summary Create a new subscription
// @Description Create a new subscription with the provided details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 201 {object} map[string]int "ID of the created subscription"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions [post]
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

func (h *SubscriptionHandler) UpdateSubscriptionsByID(c *gin.Context) {
	// TODO: implement update subscription by ID
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

// @Summary Get total price of subscriptions
// @Description Calculate the total price of subscriptions for a user within a date range, optionally filtered by service name
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID"
// @Param service_name query string false "Service Name"
// @Param start_date query string true "Start Date (MM-YYYY)"
// @Param end_date query string true "End Date (MM-YYYY)"
// @Success 200 {object} map[string]int "Total price of subscriptions"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /subscriptions/totalprice [get]
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
