// Package memory provides simple, in-memory head reader and setter for testing and development.
package memory

import (
	"context"

	"github.com/ipfs/go-cid"
)

type HeadManager struct {
	currentHead cid.Cid
}

func NewHeadManager(currentHead cid.Cid) *HeadManager {
	return &HeadManager{currentHead: currentHead}
}

func (h *HeadManager) ReadHead(ctx context.Context) (cid.Cid, error) {
	return h.currentHead, nil
}

func (h *HeadManager) SetHead(ctx context.Context, cid cid.Cid) error {
	h.currentHead = cid
	return nil
}
