package service

import (
	"context"
	"fmt"
	"log/slog"

	"gopkg.in/telebot.v4"
)

type StubSender struct {
	bot           *telebot.Bot
	log           *slog.Logger
	stubRecipient int64
}

func (s *StubSender) Send(ctx context.Context, recipient string, message string) error {
	s.log.Info("Start sending Notify",
		"recipient", recipient,
		"message", message)

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

func NewTelegramSender(log *slog.Logger, bot *telebot.Bot, stubRecipient int64) *StubSender {
	return &StubSender{
		bot:           bot,
		log:           log,
		stubRecipient: stubRecipient,
	}
}
