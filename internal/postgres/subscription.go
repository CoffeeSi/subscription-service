package postgres

import (
	"context"

	"github.com/CoffeeSi/subscription-service/internal/model"
	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository struct {
	db *Database
}

func NewSubscriptionRepository(db *Database) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (h *SubscriptionRepository) Create(sub model.Subscription) (int, error) {
	var id int
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) 
		VALUES ($1, $2, $3, TO_DATE($4, 'MM-YYYY'), TO_DATE($5, 'MM-YYYY')) RETURNING id`
	err := h.db.Connect.QueryRow(context.Background(), query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *SubscriptionRepository) List() ([]model.Subscription, error) {
	query := `SELECT service_name, price, user_id, to_char(start_date, 'MM-YYYY') AS start_date, 
		to_char(end_date, 'MM-YYYY') AS end_date FROM subscriptions`
	rows, err := h.db.Connect.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	subs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Subscription])
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (h *SubscriptionRepository) GetByID(id string) ([]model.Subscription, error) {
	query := `SELECT service_name, price, user_id, to_char(start_date, 'MM-YYYY') AS start_date, 
		to_char(end_date, 'MM-YYYY') AS end_date FROM subscriptions WHERE user_id=$1`
	rows, err := h.db.Connect.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}

	subs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Subscription])
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (h *SubscriptionRepository) DeleteByID(id string) error {
	query := "DELETE FROM subscriptions WHERE user_id=$1"
	_, err := h.db.Connect.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}

func (h *SubscriptionRepository) GetTotalPriceByID(userID string, serviceName string, startDate string, endDate string) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) 
			FROM subscriptions s
			JOIN generate_series( 
				to_date($3, 'MM-YYYY'), 
				to_date($4, 'MM-YYYY'),
				interval '1 month' 
			) AS m(month)
			ON s.start_date <= m.month
			AND (s.end_date IS NULL OR s.end_date >= m.month)
			WHERE user_id=$1 AND ($2 = '*' OR service_name=$2);`
	var total int
	err := h.db.Connect.QueryRow(context.Background(), query, userID, serviceName, startDate, endDate).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
