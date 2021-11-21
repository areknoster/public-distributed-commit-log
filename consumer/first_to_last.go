package consumer

import (
	"context"
	"github.com/areknoster/public-distributed-commit-log/head"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/rs/zerolog/log"
	"time"
)

type FirstToLastConsumerConfig struct {
	PollInterval time.Duration
	ExpirationPeriod time.Duration
}

type FirstToLastConsumer struct{
	headReader head.Reader
	consumerOffsetManager head.Manager
	storage storage.MessageStorage
	config FirstToLastConsumerConfig
}

func (f *FirstToLastConsumer) Consume(globalCtx context.Context, handler MessageHandler) error {
	wait := time.After(0) // initially don't wait at all
	for {
		select{
		case <- globalCtx.Done():
			f.syncOffset()
			return nil
		}

		wait = time.After(f.config.PollInterval)
	}
}

func (f *FirstToLastConsumer) syncOffset() {
	// if we had some caching mechanism, we would make sure that the value got persistently written here
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cid, err := f.consumerOffsetManager.ReadHead(ctxTimeout)
	if err != nil {
		log.Error().Err(err).Msg("read offset for logging before closing consumer")
	}
	log.Info().Str("offset", cid.String()).Msg("consumer closed")
}


