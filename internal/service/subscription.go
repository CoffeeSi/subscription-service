package service

import (
	"errors"
	"time"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/CoffeeSi/subscription-service/internal/repository"
	"github.com/google/uuid"
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
		if formattedStartDate.After(formattedEndDate) {
			s.logger.Error("start_date cannot be after end_date", zap.String("start_date", sub.StartDate), zap.String("end_date", *sub.EndDate))
			return 0, errors.New("start_date cannot be after end_date")
		}

		formattedEndDateStr := formattedEndDate.Format("01-2006")
		sub.EndDate = &formattedEndDateStr
	}
	if formattedStartDate.After(time.Now()) {
		return 0, errors.New("start_date cannot be in the future")
	}
	return s.subscriptions.Create(sub)
}

func (s *SubscriptionService) ListSubscriptions() ([]model.Subscription, error) {
	return s.subscriptions.List()
}

func (s *SubscriptionService) GetSubscription(id uuid.UUID) ([]model.Subscription, error) {
	return s.subscriptions.GetByID(id)
}

func (s *SubscriptionService) UpdateSubscriptionByID(id uuid.UUID, serviceName string, sub model.Subscription) error {
	formattedStartDate, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		s.logger.Error("invalid start_date format", zap.String("start_date", sub.StartDate), zap.Error(err))
		return err
	}
	sub.StartDate = formattedStartDate.Format("01-2006")

	if sub.EndDate != nil {
		formattedEndDate, err := time.Parse("01-2006", *sub.EndDate)
		if err != nil {
			s.logger.Error("invalid end_date format", zap.String("end_date", *sub.EndDate), zap.Error(err))
			return err
		}
		if formattedStartDate.After(formattedEndDate) {
			s.logger.Error("start_date cannot be after end_date", zap.String("start_date", sub.StartDate), zap.String("end_date", *sub.EndDate))
			return errors.New("start_date cannot be after end_date")
		}

		formattedEndDateStr := formattedEndDate.Format("01-2006")
		sub.EndDate = &formattedEndDateStr
	}
	if formattedStartDate.After(time.Now()) {
		return errors.New("start_date cannot be in the future")
	}
	return s.subscriptions.UpdateByID(id, serviceName, sub)
}

func (s *SubscriptionService) DeleteSubscription(id uuid.UUID) error {
	return s.subscriptions.DeleteByID(id)
}

func (s *SubscriptionService) GetTotalPriceByID(id uuid.UUID, serviceName *string, startDate string, endDate string) (int64, error) {
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
	return s.subscriptions.GetTotalPriceByID(id, serviceName, startDate, endDate)
}
