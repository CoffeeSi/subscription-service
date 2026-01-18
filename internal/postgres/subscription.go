package postgres

import (
	"context"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type SubscriptionRepository struct {
	db     *Database
	logger *zap.Logger
}

func NewSubscriptionRepository(db *Database, logger *zap.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SubscriptionRepository) Create(sub model.Subscription) (int, error) {
	var id int
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
		VALUES ($1, $2, $3, TO_DATE($4, 'MM-YYYY'), TO_DATE($5, 'MM-YYYY')) RETURNING id`
	err := r.db.Connect.QueryRow(context.Background(), query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		r.logger.Error("failed to create subscription", zap.Any("subscription", sub), zap.Error(err))
		return 0, err
	}
	return id, nil
}

func (r *SubscriptionRepository) List() ([]model.Subscription, error) {
	query := `SELECT service_name, price, user_id, to_char(start_date, 'MM-YYYY') AS start_date, 
		to_char(end_date, 'MM-YYYY') AS end_date FROM subscriptions`
	rows, err := r.db.Connect.Query(context.Background(), query)
	if err != nil {
		r.logger.Error("failed to list subscriptions", zap.Error(err))
		return nil, err
	}

	subs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Subscription])
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) ([]model.Subscription, error) {
	query := `SELECT service_name, price, user_id, to_char(start_date, 'MM-YYYY') AS start_date, 
		to_char(end_date, 'MM-YYYY') AS end_date FROM subscriptions WHERE user_id=$1`
	rows, err := r.db.Connect.Query(context.Background(), query, id)
	if err != nil {
		r.logger.Error("failed to get subscriptions by id", zap.String("user_id", id.String()), zap.Error(err))
		return nil, err
	}

	subs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Subscription])
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepository) UpdateByID(id uuid.UUID, serviceName string, sub model.Subscription) error {
	query := `UPDATE subscriptions SET service_name=$1, price=$2, start_date=TO_DATE($3, 'MM-YYYY'), 
		end_date=TO_DATE($4, 'MM-YYYY') WHERE user_id=$5 AND service_name=$6`
	_, err := r.db.Connect.Exec(context.Background(), query,
		sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, id, serviceName)
	if err != nil {
		r.logger.Error("failed to update subscriptions", zap.String("user_id", id.String()), zap.Any("subscription", sub), zap.Error(err))
		return err
	}
	return nil
}

func (r *SubscriptionRepository) DeleteByID(id uuid.UUID) error {
	query := "DELETE FROM subscriptions WHERE user_id=$1"
	_, err := r.db.Connect.Exec(context.Background(), query, id.String())
	if err != nil {
		r.logger.Error("failed to delete subscriptions", zap.String("user_id", id.String()), zap.Error(err))
		return err
	}
	return nil
}

func (r *SubscriptionRepository) GetTotalPriceByID(userID uuid.UUID, serviceName *string, startDate string, endDate string) (int64, error) {
	query := `SELECT COALESCE(SUM(price), 0) 
			FROM subscriptions s
			JOIN generate_series( 
				TO_DATE($3, 'MM-YYYY'), 
				TO_DATE($4, 'MM-YYYY'),
				INTERVAL '1 month' 
			) AS m(month)
			ON s.start_date <= m.month
			AND (s.end_date IS NULL OR s.end_date >= m.month)
			WHERE user_id=$1 AND ($2::text IS NULL OR service_name = $2::text);`
	var total int64
	err := r.db.Connect.QueryRow(context.Background(), query, userID, serviceName, startDate, endDate).Scan(&total)
	if err != nil {
		r.logger.Error("failed to get total price by id",
			zap.String("user_id", userID.String()), zap.Stringp("service_name", serviceName),
			zap.String("start_date", startDate), zap.String("end_date", endDate), zap.Error(err))
		return 0, err
	}
	return total, nil
}
