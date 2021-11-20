package storage

import (
	"errors"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"
)

type Writer interface {
	Write(message proto.Message) (cid.Cid, error)
}

// Reader gets
type Reader interface {
	Read(cid cid.Cid, message proto.Message) error
}

type Storage interface {
	Reader
	Writer
}

type Error error

var (
	ErrMarshall   = errors.New("error when marshalling message")
	ErrUnmarshall = errors.New("error when marshalling message")
	ErrInternal   = errors.New("internal storage error")
)
