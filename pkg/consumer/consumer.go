package consumer

import (
	"billing-engine/pkg/logger"
	"context"
	"github.com/IBM/sarama"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Consumer struct {
	config Config
	conn   *kafka.Conn
	log    logger.Logger

	processor MessageProcessor
}

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, payload []byte) error
}

func StartNewConsumer(cfg Config, processor MessageProcessor, log logger.Logger) error {
	saramaConfig := sarama.NewConfig()

	c, err := sarama.NewConsumer([]string{cfg.Brokers}, saramaConfig)
	if err != nil {
		log.WithField("error", err).Error("[NewConsumer] failed to create kafka consumer")
		return err
	}

	partitions, err := c.Partitions(cfg.Topic)
	if err != nil {
		log.WithField("error", err).Error("[NewConsumer] failed to get partitions")
		return err
	}

	var (
		messages = make(chan *sarama.ConsumerMessage, 256)
		closing  = make(chan struct{})
		wg       sync.WaitGroup
	)

	log.Info("[NewConsumer] starting consumer")

	for _, partition := range partitions {
		pc, err := c.ConsumePartition(cfg.Topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.WithField("error", err).Error("[NewConsumer] failed to create partition consumer")
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			<-closing
		}(pc)

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				messages <- msg
			}
		}(pc)
	}

	go func() {
		for msg := range messages {
			err := processor.ProcessMessage(context.Background(), msg.Value)
			if err != nil {
				log.WithField("error", err).WithField("consumer_name", cfg.ConsumerName).
					Error("[Consumer] failed to process message")
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	if err := c.Close(); err != nil {
		log.WithField("error", err).Error("[NewConsumer] failed to close consumer")
		return err
	}

	return nil
}
