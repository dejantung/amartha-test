package main

import (
	"billing-engine/internal/payment/repository"
	"billing-engine/internal/payment/service"
	"billing-engine/pkg/config"
	"billing-engine/pkg/database"
	"billing-engine/pkg/logger"
	"billing-engine/pkg/producer"
	"context"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
)

func main() {
	log := logger.NewZeroLogger("consumer-payment")
	cfg, err := config.NewConfig("payment")
	if err != nil {
		panic(err)
	}

	gorm, err := database.NewGormConnection(cfg)
	if err != nil {
		panic(err)
	}

	producerConfig := producer.Config{
		Brokers:      cfg.Kafka.Broker,
		Topic:        cfg.Kafka.PaymentTopic,
		WriteTimeout: cfg.Kafka.Timeout,
	}
	newProducer, err := producer.NewProducer(producerConfig, log)
	if err != nil {
		panic(err)
	}

	paymentRepository := repository.NewPaymentRepository(gorm)
	paymentService := service.NewPaymentService(paymentRepository, newProducer, log)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	brokers := []string{cfg.Kafka.Broker}

	master, err := sarama.NewConsumer(brokers, saramaConfig)
	if err != nil {
		panic(err)
	}

	consumer, err := master.ConsumePartition(cfg.Kafka.LoanTopic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.WithField("error", err).Error("consumer error")
			case msg := <-consumer.Messages():
				log.WithField("message", string(msg.Value)).Info("received message")
				err = paymentService.ProcessMessage(context.Background(), msg.Value)
				if err != nil {
					log.WithField("error", err).Error("failed to process message")
				}
			case <-signals:
				log.Info("interrupt signal received")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
}
