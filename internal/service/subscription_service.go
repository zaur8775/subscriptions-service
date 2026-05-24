package service

import (
	"context"
	"time"

	"subscriptions-service/internal/model"
	"subscriptions-service/internal/repository"
)

type SubscriptionService interface {
	Create(
		ctx context.Context,
		req *model.CreateSubscriptionRequest,
	) error
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
		req *model.UpdateSubscriptionRequest,
	) error
	GetTotalCost(
		ctx context.Context,
		userID string,
		serviceName string,
		from string,
		to string,
	) (int, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(
	repo repository.SubscriptionRepository,
) SubscriptionService {
	return &subscriptionService{
		repo: repo,
	}
}

func (s *subscriptionService) Create(
	ctx context.Context,
	req *model.CreateSubscriptionRequest,
) error {
	startDate, err := time.Parse(
		"01-2006",
		req.StartDate,
	)
	if err != nil {
		return err
	}

	var endDate *time.Time

	if req.EndDate != nil {

		parsedEndDate, err := time.Parse(
			"01-2006",
			*req.EndDate,
		)
		if err != nil {
			return err
		}

		endDate = &parsedEndDate
	}

	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return s.repo.Create(ctx, subscription)
}

func (s *subscriptionService) GetByID(
	ctx context.Context,
	id string,
) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) List(
	ctx context.Context,
) ([]model.Subscription, error) {
	return s.repo.List(ctx)
}

func (s *subscriptionService) Delete(
	ctx context.Context,
	id string,
) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) Update(
	ctx context.Context,
	id string,
	req *model.UpdateSubscriptionRequest,
) error {
	startDate, err := time.Parse(
		"01-2006",
		req.StartDate,
	)
	if err != nil {
		return err
	}

	var endDate *time.Time

	if req.EndDate != nil {

		parsedEndDate, err := time.Parse(
			"01-2006",
			*req.EndDate,
		)
		if err != nil {
			return err
		}

		endDate = &parsedEndDate
	}

	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	return s.repo.Update(
		ctx,
		id,
		subscription,
	)
}

func (s *subscriptionService) GetTotalCost(
	ctx context.Context,
	userID string,
	serviceName string,
	from string,
	to string,
) (int, error) {
	var (
		startDate time.Time
		endDate   time.Time
		err       error
	)

	startDate, err = time.Parse(
		"01-2006",
		from,
	)
	if err != nil {
		return 0, err
	}

	if to != "" {
		endDate, err = time.Parse(
			"01-2006",
			to,
		)
		if err != nil {
			return 0, err
		}
	}

	return s.repo.GetTotalCost(
		ctx,
		userID,
		serviceName,
		startDate,
		endDate,
	)
}
