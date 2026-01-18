package repository

import (
	"github.com/CoffeeSi/subscription-service/internal/model"
)

type SubscriptionRepository interface {
	Create(sub model.Subscription) (int, error)
	List() ([]model.Subscription, error)
	GetByID(id string) ([]model.Subscription, error)
	DeleteByID(id string) error
	GetTotalPriceByID(id string, serviceName *string, startDate string, endDate string) (int64, error)
}
