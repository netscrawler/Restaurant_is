package kafkarepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/netscrawler/Restaurant_is/auth/internal/domain/models"
	infraKafka "github.com/netscrawler/Restaurant_is/auth/internal/infra/out/kafka"
)

// UserEventProducer отправляет события о пользователях в Kafka
// Используйте один экземпляр на приложение.
type UserEventProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewUserEventProducer создаёт новый репозиторий для отправки событий о пользователях
func NewUserEventProducer(kafka *infraKafka.Kafka, topic string) *UserEventProducer {
	return &UserEventProducer{
		producer: kafka.Producer,
		topic:    topic,
	}
}

// UserCreatedPayload структура для события user_created
// Email или Phone может быть пустым, если не применимо
// (например, для Staff — Email, для Client — Phone)
type UserCreatedPayload struct {
	ID       string `json:"id"`
	UserType string `json:"user_type"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

// SendUserCreatedEvent отправляет событие user_created в Kafka
func (u *UserEventProducer) SendUserCreatedEvent(ctx context.Context, user any) error {
	var payload UserCreatedPayload

	switch v := user.(type) {
	case *models.Client:
		payload = UserCreatedPayload{
			ID:       v.ID.String(),
			UserType: string(models.UserTypeClient),
			Phone:    v.Phone,
		}
	case *models.Staff:
		payload = UserCreatedPayload{
			ID:       v.ID.String(),
			UserType: string(models.UserTypeStaff),
			Email:    v.WorkEmail,
		}
	default:
		return fmt.Errorf("unsupported user type: %T", user)
	}

	value, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal user_created payload: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: u.topic,
		Key:   sarama.StringEncoder("user_created"),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = u.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send user_created event: %w", err)
	}

	return nil
}
