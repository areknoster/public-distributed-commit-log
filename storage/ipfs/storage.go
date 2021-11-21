package ipfs

import (
	"context"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"
)

type Storage struct{}

func (s Storage) Read(ctx context.Context, cid cid.Cid, message proto.Message) error {
	panic("implement me")
}

func (s Storage) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	panic("implement me")
}
