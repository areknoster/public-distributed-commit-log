package storage

import (
	"context"
	"github.com/ipfs/go-cid"
)

type ContentReader interface{
	Read(ctx context.Context, cid cid.Cid) ([]byte, error)
}

type ContentWriter interface {
	Write(ctx context.Context, content []byte, cid cid.Cid)  error
}

type ContentStorage interface{
	ContentReader
	ContentWriter
}
