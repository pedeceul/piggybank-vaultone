package events

import (
	"context"
	"time"
)

type OutboxEvent struct {
	ID          string
	AggregateID string
	Type        string
	PayloadJSON []byte
	CreatedAt   time.Time
	PublishedAt *time.Time
	Attempts    int
}

type OutboxStore interface {
	Add(ctx context.Context, evt OutboxEvent) error
	MarkPublished(ctx context.Context, id string) error
	NextUnpublished(ctx context.Context, limit int) ([]OutboxEvent, error)
}
