package kafka

import (
	"log/slog"

	"github.com/IBM/sarama"
)

type Kafka struct {
	producer sarama.Cons
	topic    string
	log      *slog.Logger
}
