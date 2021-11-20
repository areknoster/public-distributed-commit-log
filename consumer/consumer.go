package consumer

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type MessageHandler interface {
	Handle(ctx context.Context, message proto.Message) error
}

type MessageFandlerFunc func(ctx context.Context, message proto.Message) error

func (m MessageFandlerFunc) Handle(ctx context.Context, message proto.Message) error {
	return m(ctx, message)
}

// Consumer defines an interface for blocking action for listening for incoming events
// and invoking handler on each of them
type Consumer interface {
	Consume(globalCtx context.Context, handler MessageHandler) error
}
