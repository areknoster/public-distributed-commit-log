package commiter

import (
	"context"
	"fmt"
	"sync"

	"github.com/benbjohnson/clock"
	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/pdclpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/thead"
)

// IntervalCommitter adds commit at given intervals.
type IntervalCommitter struct {
	headManager    thead.Manager
	messageStorage storage.MessageStorage
	pinner         sentinel.Pinner
	ipnsManager    ipns.Manager

	ticker *clock.Ticker

	mu          sync.Mutex
	uncommitted []cid.Cid
}

func NewIntervalCommitter(headManager thead.Manager, messageStorage storage.MessageStorage, pinner sentinel.Pinner,
	ipnsManager ipns.Manager, ticker *clock.Ticker) *IntervalCommitter {
	ic := &IntervalCommitter{
		headManager:    headManager,
		messageStorage: messageStorage,
		pinner:         pinner,
		ipnsManager:    ipnsManager,
		ticker:         ticker,
	}
	go ic.run()
	return ic
}

func (i *IntervalCommitter) run() {
	for {
		<-i.ticker.C
		i.mu.Lock()
		// TODO: this blocks adding new messages when committing the old ones.
		// This should be decoupled.
		if err := i.commit(); err != nil {
			log.Error().Err(err).Msg("committing messages")
			i.mu.Unlock()
			continue
		}
		i.uncommitted = []cid.Cid{}
		i.mu.Unlock()
	}
}

func (i *IntervalCommitter) commit() error {
	if len(i.uncommitted) == 0 {
		log.Debug().Msg("nothing to commit")
		return nil
	}
	commitCtx := context.Background()
	currentHead, err := i.headManager.ReadHead(commitCtx)
	if err != nil {
		return fmt.Errorf("get current head: %w", err)
	}
	var cids []string
	for _, v := range i.uncommitted {
		cids = append(cids, v.String())
	}

	commit := &pdclpb.Commit{
		Created:           timestamppb.Now(),
		PreviousCommitCid: currentHead.String(),
		MessagesCids:      cids,
	}

	commitCID, err := i.messageStorage.Write(context.Background(), commit)
	if err != nil {
		return fmt.Errorf("write message to storage: %w", err)
	}

	if err := i.pinner.Pin(commitCtx, commitCID); err != nil {
		return fmt.Errorf("pin commit: %w", err)
	}

	if err := i.headManager.SetHead(commitCtx, commitCID); err != nil {
		return fmt.Errorf("set topic head to commit cid: %w", err)
	}
	return i.ipnsManager.UpdateIPNSEntry(commitCID.String())
}

func (i *IntervalCommitter) Add(ctx context.Context, cid cid.Cid) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.uncommitted = append(i.uncommitted, cid)
	return nil
}
