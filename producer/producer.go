// Package producer defines functions for producing messages.
package producer

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// Producer defines an interface for a blocking operation of adding message to the configured pdcl topic
type Producer interface {
	Produce(ctx context.Context, message proto.Message) error
}
