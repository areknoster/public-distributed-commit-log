package storage

import (
	"context"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"
)

type (
	MessageWriter interface {
		Write(ctx context.Context, message proto.Message) (cid.Cid, error)
	}

	MessageReader interface {
		Read(ctx context.Context, cid cid.Cid) (ProtoDecodable, error)
	}

	MessageStorage interface {
		MessageReader
		MessageWriter
	}
)
