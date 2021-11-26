// Package head provides abstractions for Public Distributed Commit Log topic head manipulation.
package thead

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
)

type Reader interface {
	ReadHead(ctx context.Context) (cid.Cid, error)
}

type Setter interface {
	SetHead(ctx context.Context, cid cid.Cid) error
}

type Manager interface {
	Reader
	Setter
}

// ErrTopicNotStarted means, that there are no commits added to topic yet.
var ErrTopicNotStarted = errors.New("topic head is started")
