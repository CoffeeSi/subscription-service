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
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date) 
		VALUES ($1, $2, $3, $4) RETURNING id`
	err := h.db.Connect.QueryRow(context.Background(), query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (h *SubscriptionRepository) List() ([]model.Subscription, error) {
	query := "SELECT service_name, price, user_id, start_date FROM subscriptions"
	rows, _ := h.db.Connect.Query(context.Background(), query)
	subs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Subscription])
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (h *SubscriptionRepository) Get(id string) ([]model.Subscription, error) {
	query := "SELECT service_name, price, user_id, start_date FROM subscriptions WHERE user_id=$1"
	rows, err := h.db.Connect.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Subscription])
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (h *SubscriptionRepository) Delete(id string) error {
	query := "DELETE FROM subscriptions WHERE user_id=$1"
	_, err := h.db.Connect.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
