package storage

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/pdcl"
)

// MessageWriter persists message with CID accessor
type MessageWriter interface {
	Write(ctx context.Context, message proto.Message) (cid.Cid, error)
}

// MessageReader lets user access message based on CID address
type MessageReader interface {
	Read(ctx context.Context, cid cid.Cid) (ProtoUnmarshallable, error)
}

// ProtoUnmarshallable can be used to deserialize message data to proto structure
type ProtoUnmarshallable interface {
	Unmarshall(message proto.Message) error
}

// MessageStorage is responsible for accessing the messages based on CID
type MessageStorage interface {
	MessageReader
	MessageWriter
}

type ProtoMessageStorage struct {
	contentStorage ContentStorage
}

func NewProtoMessageStorage(contentStorage ContentStorage) *ProtoMessageStorage {
	return &ProtoMessageStorage{
		contentStorage: contentStorage,
	}
}

func (p *ProtoMessageStorage) Read(ctx context.Context, cid cid.Cid) (ProtoUnmarshallable, error) {
	content, err := p.contentStorage.Read(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("read message from content storage: %w", err)
	}
	return ProtoDecode(content), nil
}

func (p *ProtoMessageStorage) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	encoded, err := ProtoEncode(message)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("marshall message: %w", err)
	}

	messageCID, err := pdcl.CID(encoded)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("get CID from mashalled message: %s", err)
	}
	if err := p.contentStorage.Write(ctx, encoded, messageCID); err != nil {
		return cid.Cid{}, fmt.Errorf("write message to content storage: %w", err)
	}
	return messageCID, nil
}

func ProtoEncode(message proto.Message) ([]byte, error) {
	return proto.MarshalOptions{
		Deterministic: true,
	}.Marshal(message)
}

func ProtoDecode(content []byte) ProtoUnmarshallable {
	return unmarshallable{
		protoBuf: content,
		options: proto.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
}

type unmarshallable struct {
	protoBuf []byte
	options  proto.UnmarshalOptions
}

func (u unmarshallable) Unmarshall(message proto.Message) error {
	return u.options.Unmarshal(u.protoBuf, message)
}
