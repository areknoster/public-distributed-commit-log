package memory

import (
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
)

type Storage map[cid.Cid] []byte

func (s Storage) Read(ctx context.Context, cid cid.Cid) ([]byte, error) {
	content, exists := s[cid]
	if !exists{
		return nil, fmt.Errorf("doesn't exist")
	}
	return content, nil
}

func (s Storage) Write(ctx context.Context, content []byte, cid cid.Cid) error {
	s[cid] = content
	return nil
}



