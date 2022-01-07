package commiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/pdclpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/thead"
)

// MaxBufferCommitter adds commit after buffer is filled. If it's not filled for
// a long time, it adds commit after a specified interval.
type MaxBufferCommitter struct {
	headManager    thead.Manager
	messageStorage storage.MessageStorage
	pinner         sentinel.Pinner

	ticker *clock.Ticker

	mu            sync.Mutex
	uncommitted   []cid.Cid
	maxBufferSize int
}

// NewMaxBufferCommitter returns a new instance of MaxBufferCommitter.
func NewMaxBufferCommitter(headManager thead.Manager, messageStorage storage.MessageStorage, pinner sentinel.Pinner,
	ticker *clock.Ticker, maxBufferSize int) *MaxBufferCommitter {
	mbc := &MaxBufferCommitter{
		headManager:    headManager,
		messageStorage: messageStorage,
		pinner:         pinner,
		ticker:         ticker,
		maxBufferSize:  maxBufferSize,
	}
	go mbc.run()
	return mbc
}

// TODO: some of these methods are similar to IntervalCommitter and can be deduplicated.
func (mbc *MaxBufferCommitter) run() {
	for {
		<-mbc.ticker.C
		log.Debug().Msg("reached a timeout, committing messages...")
		if err := mbc.tryCommit(); err != nil {
			log.Error().Err(err).Msg("trying to commit")
		}
	}
}

func (mbc *MaxBufferCommitter) tryCommit() error {
	mbc.mu.Lock()
	defer mbc.mu.Unlock()
	if err := mbc.commit(); err != nil {
		return fmt.Errorf("committing messages: %v", err)
	}
	mbc.uncommitted = []cid.Cid{}
	return nil
}

func (mbc *MaxBufferCommitter) commit() error {
	if len(mbc.uncommitted) == 0 {
		log.Debug().Msg("nothing to commit")
		return nil
	}
	commitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	currentHead, err := mbc.headManager.ReadHead(commitCtx)
	if err != nil {
		return fmt.Errorf("get current head: %w", err)
	}
	var cids []string
	for _, v := range mbc.uncommitted {
		cids = append(cids, v.String())
	}

	commit := &pdclpb.Commit{
		Created:           timestamppb.Now(),
		PreviousCommitCid: currentHead.String(),
		MessagesCids:      cids,
	}

	commitCID, err := mbc.messageStorage.Write(context.Background(), commit)
	if err != nil {
		return fmt.Errorf("write message to storage: %w", err)
	}

	if err := mbc.pinner.Pin(commitCtx, commitCID); err != nil {
		return fmt.Errorf("pin commit: %w", err)
	}

	if err := mbc.headManager.SetHead(commitCtx, commitCID); err != nil {
		return fmt.Errorf("set topic head to commit cid: %w", err)
	}
	log.Debug().Msgf("committed %d messages", len(cids))
	return nil
}

func (mbc *MaxBufferCommitter) Add(ctx context.Context, cid cid.Cid) error {
	mbc.mu.Lock()
	defer mbc.mu.Unlock()
	mbc.uncommitted = append(mbc.uncommitted, cid)
	if len(mbc.uncommitted) >= mbc.maxBufferSize {
		go func() {
			log.Debug().Msgf("reached a maximum commit buffer of %d, committing...", mbc.maxBufferSize)
			mbc.tryCommit()
		}()
	}
	return nil
}
