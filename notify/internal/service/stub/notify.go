package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

type StubSender struct {
	bot           *telebot.Bot
	log           *zap.Logger
	stubRecipient int64
}

func (s *StubSender) Send(ctx context.Context, recipient string, message string) error {
	s.log = s.log.With(
		zap.String("recipient", recipient),
		zap.String("message", message),
	)

	s.log.Info("Start sending Notify")

	_, err := s.bot.Send(&telebot.Chat{ID: s.stubRecipient}, formatMessage(recipient, message))
	if err != nil {
		s.log.Info("Error while sending notidy")
		return err
	}
	return nil
}

func formatMessage(recipient, message string) string {
	return fmt.Sprintf("Recipient_%s_message:%s", recipient, message)
}

func NewTelegramSender(log *zap.Logger, bot *telebot.Bot, stubRecipient int64) *StubSender {
	return &StubSender{
		bot:           bot,
		log:           log,
		stubRecipient: stubRecipient,
	}
}
