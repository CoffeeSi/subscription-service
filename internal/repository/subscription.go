package repository

import "github.com/CoffeeSi/subscription-service/internal/model"

type SubscriptionRepository interface {
	Create(sub model.Subscription) (int, error)
	List() ([]model.Subscription, error)
	Get(id string) ([]model.Subscription, error)
	Delete(id string) error
}
