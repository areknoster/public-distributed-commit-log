// Package memory provides simple, in-memory head reader and setter for testing and development.
package memory

import (
	"context"

	"github.com/areknoster/public-distributed-commit-log/head"
	"github.com/ipfs/go-cid"
)

type HeadManager struct {
	currentHead cid.Cid
}

func (h *HeadManager) ReadHead(ctx context.Context) (cid.Cid, error) {
	if h.currentHead == cid.Undef {
		return cid.Cid{}, head.ErrTopicNotStarted
	}
	return h.currentHead, nil
}

func (h *HeadManager) SetHead(ctx context.Context, cid cid.Cid) error {
	h.currentHead = cid
	return nil
}
