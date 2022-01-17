// Package storage defines interfaces for content storage, message storage and codecs.
package storage

import (
	"context"

	"github.com/ipfs/go-cid"
)

// ContentReader is responsible for fetching content based on it's CID.
type ContentReader interface {
	Read(ctx context.Context, cid cid.Cid) ([]byte, error)
}

// ContentWriter writes given content with assigned CID
type ContentWriter interface {
	Write(ctx context.Context, content []byte, cid cid.Cid) error
}

// ContentStorage allows its user to access content based on CID.
type ContentStorage interface {
	ContentReader
	ContentWriter
}
