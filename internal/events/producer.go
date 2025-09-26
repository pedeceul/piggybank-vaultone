package events

import "context"

type Producer interface {
	ProduceJSON(ctx context.Context, topic string, key string, payload []byte) error
}
