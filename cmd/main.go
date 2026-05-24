package main

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "subscriptions-service/docs"
	"subscriptions-service/internal/logger"

	"subscriptions-service/internal/config"
	"subscriptions-service/internal/database"
	"subscriptions-service/internal/handler"
	"subscriptions-service/internal/repository"
	"subscriptions-service/internal/service"
)

// @title Subscription Service API
// @version 1.0
// @description REST API for managing user subscriptions
// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.MustLoad()
	db := database.New(cfg)
	defer db.Close()

	repo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(repo)
	subscriptionHandler := handler.NewSubscriptionHandler(
		subscriptionService,
	)

	mux := http.NewServeMux()

	mux.Handle(
		"GET /swagger/",
		httpSwagger.WrapHandler,
	)

	mux.HandleFunc(
		"POST /subscriptions",
		subscriptionHandler.Create,
	)

	mux.HandleFunc(
		"GET /subscriptions",
		subscriptionHandler.List,
	)

	mux.HandleFunc(
		"GET /subscriptions/total",
		subscriptionHandler.GetTotalCost,
	)

	mux.HandleFunc(
		"GET /subscriptions/{id}",
		subscriptionHandler.GetByID,
	)

	mux.HandleFunc(
		"PUT /subscriptions/{id}",
		subscriptionHandler.Update,
	)

	mux.HandleFunc(
		"DELETE /subscriptions/{id}",
		subscriptionHandler.Delete,
	)

	log.Println("server started on :" + cfg.AppPort)

	loggedMux := logger.Logging(mux)

	err := http.ListenAndServe(
		":"+cfg.AppPort,
		loggedMux,
	)
	if err != nil {
		log.Fatal(err)
	}
}
