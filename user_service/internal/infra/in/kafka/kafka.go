package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type UserCreatedPayload struct {
	ID       string `json:"id"`
	UserType string `json:"user_type"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type UserEventConsumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	handler  UserEventHandler
}

func NewUserEventConsumer(
	brokers []string,
	groupID, topic string,
	handler UserEventHandler,
) (*UserEventConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &UserEventConsumer{
		consumer: consumer,
		topic:    topic,
		handler:  handler,
	}, nil
}

func (c *UserEventConsumer) Start(ctx context.Context) error {
	handler := &userEventHandler{handler: c.handler}
	for {
		if err := c.consumer.Consume(ctx, []string{c.topic}, handler); err != nil {
			log.Printf("Error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type UserEventHandler interface {
	HandleUserCreatedEvent(ctx context.Context, id, email, phone string) error
}

type userEventHandler struct {
	handler UserEventHandler
}

func (h *userEventHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *userEventHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *userEventHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		var payload UserCreatedPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("Failed to unmarshal user_created payload: %v", err)
			continue
		}
		if err := h.handler.HandleUserCreatedEvent(sess.Context(), payload.ID, payload.Email, payload.Phone); err != nil {
			log.Printf("Failed to handle user_created event: %v", err)
			continue
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
