package sentinel

import (
	"context"

	"github.com/ipfs/go-cid"
)

// Commiter is responsible for adding commits at the head of the topic it serves.
type Commiter interface {
	Add(ctx context.Context, cid cid.Cid) error
}
