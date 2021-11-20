package ipfs

import (
	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"
)

type Storage struct{}

func (s Storage) Read(cid cid.Cid, message proto.Message) error {
	panic("implement me")
}

func (s Storage) Write(message proto.Message) (cid.Cid, error) {
	panic("implement me")
}
