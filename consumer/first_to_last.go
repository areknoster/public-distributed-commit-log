package consumer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/pdcl"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/thead"
)

type FirstToLastConsumerConfig struct {
	PollInterval time.Duration
	PollTimeout  time.Duration
	IPNSAddr     string
}

type FirstToLastConsumer struct {
	consumerOffsetManager thead.Manager // todo: this might this to be swapped to sth cached and with sync method
	commitReader          storage.MessageReader
	messageReader         storage.MessageReader
	config                FirstToLastConsumerConfig
	ipnsMgr               ipns.Resolver
}

func NewFirstToLastConsumer(consumerOffsetManager thead.Manager, commitReader, messageReader storage.MessageReader, config FirstToLastConsumerConfig,
	ipnsResolver ipns.Resolver) *FirstToLastConsumer {
	return &FirstToLastConsumer{
		consumerOffsetManager: consumerOffsetManager,
		commitReader:          commitReader,
		messageReader:         messageReader,
		config:                config,
		ipnsMgr:               ipnsResolver,
	}
}

func (f *FirstToLastConsumer) Consume(globalCtx context.Context, handler MessageHandler) error {
	pollTimer := time.NewTimer(0)
	defer f.syncOffset()
	for {
		select {
		case <-globalCtx.Done():
			return ErrContextDone
		case <-pollTimer.C:
			log.Debug().Msg("run poll")
			pollTimer.Reset(f.config.PollInterval)
			if err := f.pollWithTimeout(globalCtx, handler); err != nil {
				return fmt.Errorf("poll messages: %w ", err)
			}

		}
	}
}

func (f *FirstToLastConsumer) pollWithTimeout(globalCtx context.Context, handler MessageHandler) error {
	pollCtx, cancel := context.WithTimeout(globalCtx, f.config.PollTimeout)
	defer cancel()
	// todo: it might be good idea to add retry mechanism or other more sophisticated error handling
	if err := f.poll(pollCtx, handler); err != nil {
		return err
	}
	return nil
}

func (f *FirstToLastConsumer) poll(ctx context.Context, handler MessageHandler) error {
	currOffset, err := f.consumerOffsetManager.ReadHead(ctx)
	if err != nil {
		return fmt.Errorf("read current counsumer offset: %w", err)
	}
	topicHeadCID, err := f.ipnsMgr.ResolveIPNS(f.config.IPNSAddr)
	if err != nil {
		return fmt.Errorf("ipns resolve: %w", err)
	}
	topicHeadCID = strings.TrimPrefix(topicHeadCID, "/ipfs/")
	if topicHeadCID == "" {
		return nil // nothing new, sentinel did not set head ipns address
	}
	topicHead, err := pdcl.ParseCID(topicHeadCID)
	if err != nil {
		return fmt.Errorf("decode topic %s head: %w", topicHeadCID, err)
	}

	if currOffset == topicHead {
		return nil // nothing new, wait till new poll
	}

	handleRunner := newFirstToLastHandleRunner(f.commitReader, f.messageReader, handler, topicHead, currOffset)
	if err := handleRunner.HandleCommits(ctx); err != nil {
		return fmt.Errorf("handle commits: %w", err)
	}
	if err := f.consumerOffsetManager.SetHead(ctx, topicHead); err != nil {
		return fmt.Errorf("set new consumer offset: %w", err)
	}
	return nil
}

func (f *FirstToLastConsumer) syncOffset() {
	// todo: if we had some caching mechanism, we would make sure that the value got persistently written
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cid, err := f.consumerOffsetManager.ReadHead(ctxTimeout)
	if err != nil {
		log.Error().Err(err).Msg("read offset for logging before closing consumer")
	}
	log.Info().Str("offset", cid.String()).Msg("consumer closed")
}

type firstToLastHandleRunner struct {
	headCID        cid.Cid
	messageReader  storage.MessageReader
	commitReader   commitReader
	handler        MessageHandler
	consumerOffset cid.Cid
	// todo: it would be much less error prone if we kept index of all correctly handled messages CIDs or sth like that
}

// todo: make those configurable
const (
	defaultCommitChanLen      = 10
	defaultConcurrentHandles  = 20
	defaultReadMessageTimeout = 5 * time.Second
)

func newFirstToLastHandleRunner(
	commitMessageReader, messageReader storage.MessageReader,
	handler MessageHandler,
	headCID, consumerOffset cid.Cid) *firstToLastHandleRunner {
	commitReader := newStorageCommitReader(commitMessageReader, defaultReadMessageTimeout)

	return &firstToLastHandleRunner{
		headCID:        headCID,
		messageReader:  messageReader,
		commitReader:   commitReader,
		handler:        handler,
		consumerOffset: consumerOffset,
	}
}

func (cl *firstToLastHandleRunner) HandleCommits(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)
	commitsChan := cl.traverseCommits(ctx, group)
	messageCIDs := make(chan cid.Cid, 2*defaultConcurrentHandles)
	cl.addWorkers(ctx, group, messageCIDs)
	group.Go(func() error {
		defer close(messageCIDs)
		var lastCommit commit
		for lastCommit = range commitsChan {
			for _, messageCid := range lastCommit.Messages {
				select {
				case <-ctx.Done():
					return fmt.Errorf("did not handle all messages: %w", ErrContextDone)
				default:
					messageCIDs <- messageCid
				}
			}
		}
		return nil
	})
	return group.Wait()
}

// traverse commits makes sure, that
func (cl *firstToLastHandleRunner) traverseCommits(ctx context.Context, group *errgroup.Group) <-chan commit {
	currentCommit := commit{ // this lets us use logic below for the head commit too.
		Previous: cl.headCID,
	}

	commitsChan := make(chan commit, defaultCommitChanLen)

	group.Go(func() error {
		defer func() {
			close(commitsChan)
			log.Debug().Msg("finished traversing commits")
		}()
		// todo: add message expiration mechanism
		for {
			logCtx := log.With().Str("last_visited_commit", currentCommit.Cid.String()).Logger()
			if currentCommit.Previous == cl.consumerOffset {
				logCtx.Debug().Msg("all commits up to the consumer offset were traversed")
				return nil
			}
			if currentCommit.Previous == cid.Undef {
				return nil
			}
			commit, err := cl.commitReader.GetCommit(ctx, currentCommit.Previous)
			if err != nil {
				return fmt.Errorf("get previous message: %w", err)
			}
			currentCommit = commit
			commitsChan <- commit
		}
	})

	return commitsChan
}

func (cl *firstToLastHandleRunner) addWorkers(ctx context.Context, group *errgroup.Group, messageCIDs <-chan cid.Cid) {
	for i := 0; i < defaultConcurrentHandles; i++ {
		group.Go(func() error {
			for messageCID := range messageCIDs {
				unmarshallable, err := cl.messageReader.Read(ctx, messageCID)
				if err != nil {
					return fmt.Errorf("can't read message: %w", err)
				}
				if err := cl.handler.Handle(ctx, unmarshallable); err != nil {
					return fmt.Errorf("error when handling message: %w", err)
				}
			}
			return nil
		})
	}
}
