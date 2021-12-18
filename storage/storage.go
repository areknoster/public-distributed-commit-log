package storage

import (
	"context"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"
)

// MessageWriter persists message with CID accessor
type MessageWriter interface {
	Write(ctx context.Context, message proto.Message) (cid.Cid, error)
}

// MessageReader lets user access message based on CID address
type MessageReader interface {
	Read(ctx context.Context, cid cid.Cid) (ProtoUnmarshallable, error)
}

// ProtoUnmarshallable can be used to deserialize message data to proto structure
// todo: refactor this to just passing generic type when go 1.18 is out
type ProtoUnmarshallable interface {
	Unmarshall(message proto.Message) error
}

// MessageStorage is responsible for accessing the messages based on CID
type MessageStorage interface {
	MessageReader
	MessageWriter
}
