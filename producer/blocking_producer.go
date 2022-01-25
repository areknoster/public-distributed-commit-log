package producer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
)

type BlockingProducer struct {
	storage        storage.MessageWriter
	sentinelClient sentinelpb.SentinelClient
}

func NewBlockingProducer(writer storage.MessageWriter, sentinelClient sentinelpb.SentinelClient) *BlockingProducer {
	return &BlockingProducer{storage: writer, sentinelClient: sentinelClient}
}

func (m *BlockingProducer) Produce(ctx context.Context, message proto.Message) error {
	cid, err := m.storage.Write(ctx, message)
	if err != nil {
		return fmt.Errorf("save message to storage: %w", err)
	}
	log.Debug().Stringer("cid", cid).Msg("message stored")

	_, err = m.sentinelClient.Publish(ctx, &sentinelpb.PublishRequest{Cid: cid.String()})
	if err != nil {
		return fmt.Errorf("publish MessageBuf to sentinel: %w", err)
	}
	log.Info().Stringer("cid", cid).Msg("message accepted by sentinel")
	return nil
}
