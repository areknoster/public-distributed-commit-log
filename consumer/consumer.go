package consumer

import (
	"context"

	"github.com/areknoster/public-distributed-commit-log/storage"
)

type MessageHandler interface {
	Handle(ctx context.Context, message storage.ProtoUnmarshallable) error
}

type MessageFandlerFunc func(ctx context.Context, message storage.ProtoUnmarshallable) error

func (m MessageFandlerFunc) Handle(ctx context.Context, message storage.ProtoUnmarshallable) error {
	return m(ctx, message)
}

// Consumer defines an interface for blocking action for listening for incoming events
// and invoking handler on each of them
type Consumer interface {
	Consume(globalCtx context.Context, handler MessageHandler) error
}
