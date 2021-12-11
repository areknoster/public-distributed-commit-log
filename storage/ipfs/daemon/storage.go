package ipfs

import (
	"context"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/storage"
)

// Storage is IPFS-based storage.Storage interface implementation
type Storage struct{}

// Read tries to find message in IPFS
func (s *Storage) Read(ctx context.Context, cid cid.Cid) (storage.ProtoUnmarshallable, error) {
	panic("implement me")
}

// Write writes message to IPFS
func (s *Storage) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	panic("implement me")
}
