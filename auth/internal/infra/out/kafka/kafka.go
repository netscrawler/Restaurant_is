package kafka

import (
	"github.com/IBM/sarama"
	"github.com/netscrawler/Restaurant_is/auth/internal/config"
)

type Kafka struct {
	Producer sarama.SyncProducer
}

// NewKafka создает новый экземпляр Kafka с синхронным producer.
func NewKafka(cfg config.Kafka) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}
	return &Kafka{Producer: producer}, nil
}
