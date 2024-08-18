package server

import (
	"billing-engine/internal/billing/api"
	"billing-engine/internal/billing/repository"
	"billing-engine/internal/billing/service"
	"billing-engine/pkg/config"
	"billing-engine/pkg/database"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/producer"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Echo *echo.Echo
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewServer() (*Server, error) {
	cfg, err := config.NewConfig("billing")
	if err != nil {
		return nil, err
	}

	log := logger.NewZeroLogger("billing")
	gorm, err := database.NewGormConnection(cfg)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Cache.Host,
		DB:   cfg.Cache.Database,
	})

	newKafkaConfig := producer.Config{
		Brokers:      cfg.Kafka.Broker,
		Topic:        cfg.Kafka.LoanTopic,
		WriteTimeout: cfg.Kafka.Timeout,
	}

	kafkaProducer, err := producer.NewProducer(newKafkaConfig, log)
	if err != nil {
		return nil, err
	}

	newBillingRepository := repository.NewBillingRepositoryProvider(gorm, log)
	newBillingCache := repository.NewBillingCacheProvider(redisClient)
	billingService := service.NewBillingService(newBillingRepository, newBillingCache, kafkaProducer, log)
	billingHandler := api.NewBillingHandler(billingService)

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	billingHandler.AddRoutes(e)

	return &Server{
		Echo: e,
	}, nil
}
