package bus

import (
	"context"
)

type EventBus interface {
	Publish(ctx context.Context, event interface{}) error
	Subscribe(eventType string, handler EventHandler) error
	Start(ctx context.Context) error
	Close() error
}

type EventHandler func(ctx context.Context, event interface{}) error