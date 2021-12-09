// Package memory provides simple, in-memory head reader and setter for testing and development.
package memory

import (
	"context"
	"sync"

	"github.com/ipfs/go-cid"
)

type HeadManager struct {
	mx          sync.RWMutex
	currentHead cid.Cid
}

func NewHeadManager(currentHead cid.Cid) *HeadManager {
	return &HeadManager{currentHead: currentHead}
}

func (h *HeadManager) ReadHead(ctx context.Context) (cid.Cid, error) {
	h.mx.RLock()
	c := h.currentHead
	h.mx.RUnlock()
	return c, nil
}

func (h *HeadManager) SetHead(ctx context.Context, cid cid.Cid) error {
	h.mx.Lock()
	h.currentHead = cid
	h.mx.Unlock()
	return nil
}
