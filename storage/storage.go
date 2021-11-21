package storage

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"
)

type MessageWriter interface {
	Write(ctx context.Context, message proto.Message) (cid.Cid, error)
}

type MessageReader interface {
	Read(ctx context.Context, cid cid.Cid, message proto.Message) error
}

type MessageStorage interface {
	MessageReader
	MessageWriter
}

type Error error

var (
	ErrMarshall   = errors.New("error when marshalling message")
	ErrUnmarshall = errors.New("error when unmarshalling message")
	ErrInternal   = errors.New("internal storage error")
)
