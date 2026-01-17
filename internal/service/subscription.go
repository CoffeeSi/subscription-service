package service

import (
	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/CoffeeSi/subscription-service/internal/repository"
)

type Subscription struct {
	subscriptions repository.SubscriptionRepository
}

func NewSubscription(subscriptions repository.SubscriptionRepository) *Subscription {
	return &Subscription{
		subscriptions: subscriptions,
	}
}

func (s *Subscription) CreateSubscription(sub model.Subscription) (int, error) {
	return s.subscriptions.Create(sub)
}

func (s *Subscription) ListSubscriptions() ([]model.Subscription, error) {
	return s.subscriptions.List()
}

func (s *Subscription) GetSubscription(id string) ([]model.Subscription, error) {
	return s.subscriptions.Get(id)
}

func (s *Subscription) DeleteSubscription(id string) error {
	return s.subscriptions.Delete(id)
}
