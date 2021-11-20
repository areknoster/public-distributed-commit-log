package sentinel

import (
	"context"

	"github.com/ipfs/go-cid"
)

type Pinner interface {
	Pin(ctx context.Context, cid cid.Cid) error
}
