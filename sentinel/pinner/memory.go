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

// NewMemoryPinner initializes memory pinner which implements sentinel.Pinner interface.
func NewMemoryPinner() *MemoryPinner {
	return &MemoryPinner{
		pins: map[cid.Cid]struct{}{},
	}
}

// Pin adds cid to local pins index. It might be used to debug pinning logic.
func (l *MemoryPinner) Pin(ctx context.Context, cid cid.Cid) error {
	l.pins[cid] = struct{}{}
	return nil
}
