package handler

import (
	"encoding/json"
	"net/http"

	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(
	service service.SubscriptionService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

// Create godoc
//
// @Summary Create subscription
// @Description Create new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body model.CreateSubscriptionRequest true "Subscription data"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req model.CreateSubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Create(r.Context(), &req)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetByID godoc
//
// @Summary Get subscription by ID
// @Description Get subscription by UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")

	subscription, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			"subscription not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(subscription)
	if err != nil {

		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)

		return
	}
}

// List godoc
//
// @Summary List subscriptions
// @Description Get all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	subscriptions, err := h.service.List(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to get subscriptions",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(subscriptions)
	if err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}

// Delete godoc
//
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 404
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")

	err := h.service.Delete(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			"failed to delete subscription",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Update godoc
//
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Param id path string true "Subscription ID"
// @Param request body model.UpdateSubscriptionRequest true "Updated data"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")

	var req model.UpdateSubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Update(
		r.Context(),
		id,
		&req,
	)
	if err != nil {
		http.Error(
			w,
			"failed to update subscription",
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetTotalCost godoc
//
// @Summary Calculate total subscription cost
// @Description Calculate total cost for period
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID"
// @Param service_name query string false "Service name"
// @Param from query string true "From MM-YYYY"
// @Param to query string false "To MM-YYYY"
// @Success 200 {object} model.TotalResponse
// @Failure 500
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) GetTotalCost(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	total, err := h.service.GetTotalCost(
		r.Context(),
		userID,
		serviceName,
		from,
		to,
	)
	if err != nil {
		http.Error(
			w,
			"failed to calculate total",
			http.StatusInternalServerError,
		)
		return
	}

	response := model.TotalResponse{
		Total: total,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}
