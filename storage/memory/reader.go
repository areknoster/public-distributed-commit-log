package memory

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
)

// Storage is in-memory implementation of message story. Should be used only for testing
type Storage map[cid.Cid][]byte

// Read gets message form im-memory map
func (s Storage) Read(ctx context.Context, cid cid.Cid) ([]byte, error) {
	content, exists := s[cid]
	if !exists {
		return nil, fmt.Errorf("doesn't exist")
	}
	return content, nil
}

// Read writes message to in-memory map
func (s Storage) Write(ctx context.Context, content []byte, cid cid.Cid) error {
	s[cid] = content
	return nil
}
