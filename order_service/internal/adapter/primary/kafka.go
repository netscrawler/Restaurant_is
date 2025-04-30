package primary

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/netscrawler/Restaurant_is/order_service/internal/config"
	dto "github.com/netscrawler/Restaurant_is/order_service/internal/models/repository"
)

type KafkaPublisher struct {
	producer sarama.SyncProducer
	topic    string
	log      *slog.Logger
}

func NewKafkaPublisher(cfg *config.Kafka, log *slog.Logger) (*KafkaPublisher, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = cfg.RetryMax
	config.Producer.Return.Successes = cfg.ReturnSuccesses

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaPublisher{
		producer: producer,
		topic:    cfg.Topic,
		log:      log.With("component", "kafka_publisher"),
	}, nil
}

func (k *KafkaPublisher) PublishEvent(ctx context.Context, event *dto.Event) error {
	logger := k.log.With("func", "PublishEvent")

	payload, err := json.Marshal(event)
	if err != nil {
		logger.ErrorContext(ctx, "failed to marshal event",
			slog.String("error", err.Error()),
			slog.String("event_id", event.ID))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Key:   sarama.StringEncoder(event.ID),
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		logger.ErrorContext(ctx, "failed to send message",
			slog.String("error", err.Error()),
			slog.String("event_id", event.ID))
		return err
	}

	logger.InfoContext(ctx, "message sent successfully",
		slog.String("event_id", event.ID),
		slog.Int64("partition", int64(partition)),
		slog.Int64("offset", offset))

	return nil
}

func (k *KafkaPublisher) Close() error {
	return k.producer.Close()
}
