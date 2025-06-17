package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/netscrawler/Restaurant_is/order_service/internal/config"
	"github.com/netscrawler/Restaurant_is/order_service/internal/models/repository"
)

type EventRepository interface {
	Save(ctx context.Context, event *repository.Event) error
	GetUnpublishedEvents(ctx context.Context) ([]*repository.Event, error)
	MarkAsPublished(ctx context.Context, eventID string) error
}

type KafkaPublisher interface {
	PublishEvent(ctx context.Context, event *repository.Event) error
}

type Event struct {
	repo           EventRepository
	publisher      KafkaPublisher
	processTimeout time.Duration
	log            *slog.Logger
}

func NewEventService(
	repo EventRepository,
	publisher KafkaPublisher,
	cfg *config.Config,
	log *slog.Logger,
) *Event {
	l := log.With(slog.String("worker", "event"))
	return &Event{
		repo:           repo,
		publisher:      publisher,
		processTimeout: cfg.ProcessTimeout,
		log:            l,
	}
}

func (s *Event) SaveEvent(ctx context.Context, event *repository.Event) error {
	return s.repo.Save(ctx, event)
}

// StartBackgroundProcessing starts asynchronous processing of unpublished events
func (s *Event) StartBackgroundProcessing(ctx context.Context) {
	s.log.Info("starting background", slog.Any("process timeout", s.processTimeout))
	go func() {
		ticker := time.NewTicker(s.processTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.PublishUnpublishedEvents(ctx); err != nil {
					// Here you might want to add logging
					continue
				}
			default:
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()
}

func (s *Event) PublishUnpublishedEvents(ctx context.Context) error {
	events, err := s.repo.GetUnpublishedEvents(ctx)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := s.publisher.PublishEvent(ctx, event); err != nil {
			return err
		}

		if err := s.repo.MarkAsPublished(ctx, event.ID); err != nil {
			return err
		}
	}

	return nil
}
