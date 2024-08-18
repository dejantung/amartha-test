package server

import (
	"billing-engine/internal/billing/api"
	"billing-engine/internal/billing/domain"
	"billing-engine/internal/billing/repository"
	"billing-engine/internal/billing/service"
	"billing-engine/pkg/config"
	"billing-engine/pkg/database"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/producer"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	Echo *echo.Echo
	Log  logger.Logger
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewServer(log logger.Logger, cfg *config.Config) (*Server, error) {
	gorm, err := database.NewGormConnection(cfg)
	if err != nil {
		return nil, err
	}

	err = gorm.AutoMigrate(&domain.Customer{}, &domain.Loan{}, &domain.Schedule{})
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

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	billingHandler.AddRoutes(e)
	e.Validator = &CustomValidator{validator: validator.New()}

	return &Server{
		Echo: e,
		Log:  log,
	}, nil
}

func (s *Server) Stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	s.Log.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Echo.Shutdown(ctx); err != nil {
		s.Echo.Logger.Fatal(err)
	}

	s.Log.Info("Server shutdown gracefully")
}
