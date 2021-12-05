package consumer

import (
	"context"

	"github.com/areknoster/public-distributed-commit-log/storage"
)

// MessageHandler is called to handle message found by Consumer.
// Handle on multiple messages can be called concurrently.
type MessageHandler interface {
	Handle(ctx context.Context, message storage.ProtoUnmarshallable) error
}

// MessageHandlerFunc is a function implementing MessageHandler interface
type MessageHandlerFunc func(ctx context.Context, message storage.ProtoUnmarshallable) error

// Handle calls MessageHandlerFunc
func (m MessageHandlerFunc) Handle(ctx context.Context, message storage.ProtoUnmarshallable) error {
	return m(ctx, message)
}

// Consumer defines an interface for blocking action for listening for incoming events
// and invoking handler on each of them
type Consumer interface {
	Consume(globalCtx context.Context, handler MessageHandler) error
}
