package kafka

import (
	"billing-engine/pkg/logger"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_producer.go -package=mocks billing-engine/pkg/kafka ProducerProvider
type ProducerProvider interface {
	SendMessage(ctx context.Context, payload Message) error
}

type impl struct {
	config Config
	conn   *kafka.Conn
	log    logger.Logger
}

func (i impl) SendMessage(ctx context.Context, payload Message) error {
	msgBytes, err := json.Marshal(payload)
	if err != nil {
		i.log.WithField("error", err).
			WithField("payload", payload).Error("[SendMessage] failed to marshal payload")

		return err
	}

	newKafkaMessage := kafka.Message{
		Key:   []byte(payload.EventID),
		Value: msgBytes,
	}

	i.log.WithField("payload", payload).Info("[SendMessage] sending message to kafka")
	err = i.conn.SetWriteDeadline(time.Now().Add(time.Duration(i.config.WriteTimeout) * time.Second))
	if err != nil {
		i.log.WithField("error", err).Error("[SendMessage] failed to set write deadline")
		return err
	}

	_, err = i.conn.WriteMessages(newKafkaMessage)
	if err != nil {
		i.log.WithField("error", err).Error("[SendMessage] failed to write message")
		return err
	}

	i.log.WithField("payload", payload).Info("[SendMessage] message sent to kafka")
	return nil
}

func NewProducer(config Config, log logger.Logger) (ProducerProvider, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", config.Brokers, config.Topic, config.Partition)
	if err != nil {
		return nil, err
	}

	return &impl{
		config: config,
		conn:   conn,
		log:    log,
	}, nil
}
