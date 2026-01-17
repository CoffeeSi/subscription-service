package service

import (
	"errors"
	"time"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/CoffeeSi/subscription-service/internal/repository"
)

type SubscriptionService struct {
	subscriptions repository.SubscriptionRepository
}

func NewSubscriptionService(subscriptions repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptions: subscriptions,
	}
}

func (s *SubscriptionService) CreateSubscription(sub model.Subscription) (int, error) {
	return s.subscriptions.Create(sub)
}

func (s *SubscriptionService) ListSubscriptions() ([]model.Subscription, error) {
	return s.subscriptions.List()
}

func (s *SubscriptionService) GetSubscription(id string) ([]model.Subscription, error) {
	return s.subscriptions.GetByID(id)
}

func (s *SubscriptionService) DeleteSubscription(id string) error {
	return s.subscriptions.DeleteByID(id)
}

func (s *SubscriptionService) GetTotalPriceByID(userID string, serviceName string, startDate string, endDate string) (int, error) {
	if serviceName == "" {
		serviceName = "*"
	}
	startDateParsed, err := time.Parse("01-2006", startDate)
	if err != nil {
		return 0, err
	}
	endDateParsed, err := time.Parse("01-2006", endDate)
	if err != nil {
		return 0, err
	}
	if startDateParsed.After(endDateParsed) {
		return 0, errors.New("startDate cannot be after endDate")
	}
	return s.subscriptions.GetTotalPriceByID(userID, serviceName, startDate, endDate)
}
