package producer

import (
	"context"
	"fmt"

	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"google.golang.org/protobuf/proto"
)

type MessageProducer struct {
	storage        storage.MessageStorage
	sentinelClient sentinelpb.SentinelClient
}

func NewMessageProducer(storage storage.MessageStorage, sentinelClient sentinelpb.SentinelClient) *MessageProducer {
	return &MessageProducer{storage: storage, sentinelClient: sentinelClient}
}

func (m *MessageProducer) Produce(ctx context.Context, message proto.Message) error {
	cid, err := m.storage.Write(ctx, message)
	if err != nil {
		return fmt.Errorf("save message to storage: %w", err)
	}

	_, err = m.sentinelClient.Publish(ctx, &sentinelpb.PublishRequest{Cid: cid.String()})
	if err != nil {
		return fmt.Errorf("publish message to sentinel: %w", err)
	}
	return nil
}
