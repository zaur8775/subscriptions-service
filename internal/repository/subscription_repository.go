package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"subscriptions-service/internal/model"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *model.Subscription) error
	GetByID(
		ctx context.Context,
		id string,
	) (*model.Subscription, error)
	List(
		ctx context.Context,
	) ([]model.Subscription, error)
	Delete(
		ctx context.Context,
		id string,
	) error
	Update(
		ctx context.Context,
		id string,
		subscription *model.Subscription,
	) error
	GetTotalCost(
		ctx context.Context,
		userID string,
		serviceName string,
		startDate time.Time,
		endDate time.Time,
	) (int, error)
}

type subscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

func (r *subscriptionRepository) Create(
	ctx context.Context,
	subscription *model.Subscription,
) error {
	query := `
		INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	)

	return err
}

func (r *subscriptionRepository) GetByID(
	ctx context.Context,
	id string,
) (*model.Subscription, error) {
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var subscription model.Subscription

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *subscriptionRepository) List(
	ctx context.Context,
) ([]model.Subscription, error) {
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subscriptions []model.Subscription

	for rows.Next() {

		var subscription model.Subscription

		err := rows.Scan(
			&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		subscriptions = append(
			subscriptions,
			subscription,
		)
	}

	return subscriptions, nil
}

func (r *subscriptionRepository) Delete(
	ctx context.Context,
	id string,
) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}

func (r *subscriptionRepository) Update(
	ctx context.Context,
	id string,
	subscription *model.Subscription,
) error {
	query := `
		UPDATE subscriptions
		SET
			service_name = $1,
			price = $2,
			start_date = $3,
			end_date = $4,
			updated_at = NOW()
		WHERE id = $5
	`

	_, err := r.db.Exec(
		ctx,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.StartDate,
		subscription.EndDate,
		id,
	)

	return err
}

func (r *subscriptionRepository) GetTotalCost(
	ctx context.Context,
	userID string,
	serviceName string,
	startDate time.Time,
	endDate time.Time,
) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE user_id = $1
		AND ($2 = '' OR service_name = $2)
		AND (
			$4::timestamp = '0001-01-01 00:00:00'
			OR start_date <= $4
		)
		AND (
			end_date IS NULL
			OR end_date >= $3
		)
	`

	var total int

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		serviceName,
		endDate,
		startDate,
	).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
