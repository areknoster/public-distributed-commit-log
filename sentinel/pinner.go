package sentinel

import (
	"context"

	"github.com/ipfs/go-cid"
)

// Pinner is responsible for pinning message - thus replicating it and making accessible to external consumers
type Pinner interface {
	Pin(ctx context.Context, cid cid.Cid) error
}
