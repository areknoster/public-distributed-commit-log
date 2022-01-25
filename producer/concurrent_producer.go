package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ConcurrentProducer can run multiple Produce jobs concurrently.
// error channel should be sunk, otherwise it would block messages production eventually
// in order to stop it gracefully, close Messages channel
// and then sink Errors channel till it's closed
type ConcurrentProducer interface {
	Messages() chan<- proto.Message
	Errors() <-chan Error
}

type Error struct {
	Err     error
	Message proto.Message
}

func (e Error) Error() string {
	return fmt.Sprintf("message: \n%s\n could not be produced: %s",
		protojson.Format(e.Message),
		e.Err.Error())
}

func (e Error) Unwrap() error {
	return e.Err
}

type BasicConcurrentProducerConfig struct {
	// JobsNumber should be relatively high, since IPFS communication can sometimes block for extend periods
	JobsNumber     uint          `envconfig:"JOBS_NUMBER" default:"20"`
	ProduceTimeout time.Duration `envconfig:"PRODUCE_TIMEOUT" default:"2m"`
	ErrBuf         uint          `envconfig:"ERR_CHAN_BUFF" default:"50"`
	MessageBuf     uint          `envconfig:"MESSAGE_BUFF" default:"250"`
}

func StartBasicConcurrentProducer(globalCtx context.Context, blockingProducer Producer, config BasicConcurrentProducerConfig) *BasicConcurrentProducer {
	concurrentProducer := &BasicConcurrentProducer{
		blockingProducer: blockingProducer,
		config:           config,
		messages:         make(chan proto.Message, config.MessageBuf),
		errors:           make(chan Error, config.ErrBuf),
	}
	concurrentProducer.jobsWg.Add(int(config.JobsNumber))
	for i := uint(0); i < config.JobsNumber; i++ {
		go concurrentProducer.Job(globalCtx, i)
	}

	go func() {
		concurrentProducer.jobsWg.Wait()
		close(concurrentProducer.errors)
		log.Info().Msg("concurrent producer errors chan closed")
	}()
	return concurrentProducer
}

type BasicConcurrentProducer struct {
	blockingProducer Producer
	config           BasicConcurrentProducerConfig
	messages         chan proto.Message
	errors           chan Error
	jobsWg           sync.WaitGroup
}

func (b *BasicConcurrentProducer) Job(globalCtx context.Context, index uint) {
	for m := range b.messages {
		ctx, cancel := context.WithTimeout(globalCtx, b.config.ProduceTimeout)
		if err := b.blockingProducer.Produce(ctx, m); err != nil {
			b.errors <- Error{
				Err:     err,
				Message: m,
			}
		}
		cancel()
	}
	b.jobsWg.Done()
	log.Info().Uint("producer_index", index).Msg("producer job finished")
}

func (b *BasicConcurrentProducer) Messages() chan<- proto.Message {
	return b.messages
}

func (b *BasicConcurrentProducer) Errors() <-chan Error {
	return b.errors
}
