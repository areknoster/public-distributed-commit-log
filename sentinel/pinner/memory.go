package pinner

import (
	"context"

	"github.com/ipfs/go-cid"
)

// MemoryPinner adds a message cid to its memory.
// It can be used only in local testing, when localfs implementation is used for message storage
type MemoryPinner struct {
	pins map[cid.Cid]struct{}
}

func NewMemoryPinner() *MemoryPinner {
	return &MemoryPinner{
		pins: map[cid.Cid]struct{}{},
	}
}

func (l *MemoryPinner) Pin(ctx context.Context, cid cid.Cid) error {
	l.pins[cid] = struct{}{}
	return nil
}
