package repository

import (
	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(sub model.Subscription) (int, error)
	List() ([]model.Subscription, error)
	GetByID(id uuid.UUID) ([]model.Subscription, error)
	UpdateByID(id uuid.UUID, serviceName string, sub model.Subscription) error
	DeleteByID(id uuid.UUID) error
	GetTotalPriceByID(id uuid.UUID, serviceName *string, startDate string, endDate string) (int64, error)
}
