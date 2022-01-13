package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ipfs/go-cid"
)

func NewStorage() *Storage {
	return new(Storage)
}

// Storage is in-memory implementation of content storage. Should be used only for testing
type Storage struct {
	smap sync.Map
}

// Read gets message form im-memory map
func (s *Storage) Read(ctx context.Context, cid cid.Cid) ([]byte, error) {
	content, exists := s.smap.Load(cid)
	if !exists {
		return nil, fmt.Errorf("doesn't exist")
	}
	return content.([]byte), nil
}

// Read writes message to in-memory map
func (s *Storage) Write(ctx context.Context, content []byte, cid cid.Cid) error {
	s.smap.Store(cid, content)
	return nil
}
