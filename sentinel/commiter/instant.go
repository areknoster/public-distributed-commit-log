// Package commiter contains implementations of committers.
package commiter

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/pdclpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/thead"
)

// Instant Commiter adds commit after every single received commit
type Instant struct {
	headManager    thead.Manager
	messageStorage storage.MessageStorage
	pinner         sentinel.Pinner
	ipnsManager    ipns.Manager
}

func NewInstant(headManager thead.Manager, messageStorage storage.MessageStorage, pinner sentinel.Pinner,
	ipnsManager ipns.Manager) *Instant {
	return &Instant{
		headManager:    headManager,
		messageStorage: messageStorage,
		pinner:         pinner,
		ipnsManager:    ipnsManager,
	}
}

func (i *Instant) Add(ctx context.Context, cid cid.Cid) error {
	currentHead, err := i.headManager.ReadHead(ctx)
	if err != nil {
		return fmt.Errorf("get current head: %w", err)
	}

	commit := &pdclpb.Commit{
		Created:           timestamppb.Now(),
		PreviousCommitCid: currentHead.String(),
		MessagesCids:      []string{cid.String()},
	}

	commitCID, err := i.messageStorage.Write(nil, commit)
	if err != nil {
		return fmt.Errorf("write message to storage: %w", err)
	}

	if err := i.pinner.Pin(ctx, commitCID); err != nil {
		return fmt.Errorf("pin commit: %w", err)
	}

	if err := i.headManager.SetHead(ctx, commitCID); err != nil {
		return fmt.Errorf("set topic head to commit cid: %w", err)
	}
	return i.ipnsManager.UpdateIPNSEntry(commitCID.String())
}
