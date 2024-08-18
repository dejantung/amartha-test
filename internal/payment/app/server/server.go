package server

import (
	"billing-engine/internal/payment/api"
	"billing-engine/internal/payment/domain"
	"billing-engine/internal/payment/repository"
	"billing-engine/internal/payment/service"
	"billing-engine/pkg/config"
	"billing-engine/pkg/database"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/producer"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	err = gorm.AutoMigrate(&domain.Loan{}, &domain.PaymentSchedule{}, &domain.Payment{})
	if err != nil {
		return nil, err
	}

	producerConfig := producer.Config{
		Brokers:      cfg.Kafka.Broker,
		Topic:        cfg.Kafka.PaymentTopic,
		WriteTimeout: cfg.Kafka.Timeout,
	}
	newProducer, err := producer.NewProducer(producerConfig, log)
	if err != nil {
		return nil, err
	}

	paymentRepository := repository.NewPaymentRepository(gorm)
	paymentService := service.NewPaymentService(paymentRepository, newProducer, log)
	paymentHandler := api.NewPaymentHandler(paymentService, log)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	paymentHandler.AddRoutes(e)

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
