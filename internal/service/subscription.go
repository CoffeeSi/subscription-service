package service

import (
	"errors"
	"time"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/CoffeeSi/subscription-service/internal/repository"
	"go.uber.org/zap"
)

type SubscriptionService struct {
	logger        *zap.Logger
	subscriptions repository.SubscriptionRepository
}

func NewSubscriptionService(subscriptions repository.SubscriptionRepository, logger *zap.Logger) *SubscriptionService {
	return &SubscriptionService{
		subscriptions: subscriptions,
		logger:        logger,
	}
}

func (s *SubscriptionService) CreateSubscription(sub model.Subscription) (int, error) {
	formattedStartDate, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		s.logger.Error("invalid start_date format", zap.String("start_date", sub.StartDate), zap.Error(err))
		return 0, err
	}
	sub.StartDate = formattedStartDate.Format("01-2006")

	if sub.EndDate != nil {
		formattedEndDate, err := time.Parse("01-2006", *sub.EndDate)
		if err != nil {
			s.logger.Error("invalid end_date format", zap.String("end_date", *sub.EndDate), zap.Error(err))
			return 0, err
		}
		formattedEndDateStr := formattedEndDate.Format("01-2006")
		sub.EndDate = &formattedEndDateStr
	}
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

func (s *SubscriptionService) GetTotalPriceByID(userID string, serviceName *string, startDate string, endDate string) (int64, error) {
	startDateParsed, err := time.Parse("01-2006", startDate)
	if err != nil {
		s.logger.Error("invalid start_date format", zap.String("start_date", startDate), zap.Error(err))
		return 0, err
	}
	endDateParsed, err := time.Parse("01-2006", endDate)
	if err != nil {
		s.logger.Error("invalid end_date format", zap.String("end_date", endDate), zap.Error(err))
		return 0, err
	}
	if startDateParsed.After(endDateParsed) {
		s.logger.Warn("start_date cannot be after end_date", zap.String("start_date", startDate), zap.String("end_date", endDate))
		return 0, errors.New("start_date cannot be after end_date")
	}
	return s.subscriptions.GetTotalPriceByID(userID, serviceName, startDate, endDate)
}
