package storage

import (
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"google.golang.org/protobuf/proto"
)

type MessageWriter interface {
	Write(ctx context.Context, message proto.Message )  (cid.Cid, error)
}

type MessageReader interface {
	Read(ctx context.Context, cid cid.Cid) (ProtoUnmarshallable, error)
}

// ProtoUnmarshallable can be used to deserialize message data to proto structure
type ProtoUnmarshallable interface {
	Unmarshall(message proto.Message) error
}


type MessageStorage interface {
	MessageReader
	MessageWriter
}


type ProtoMessageStorage struct{
	contentStorage ContentStorage
	marshalOpts     proto.MarshalOptions
}

func NewProtoMessageStorage(contentStorage ContentStorage) *ProtoMessageStorage {
	return &ProtoMessageStorage{
		contentStorage: contentStorage,
		marshalOpts: proto.MarshalOptions{
		Deterministic:     true,
	}}
}

func (p *ProtoMessageStorage) Read(ctx context.Context, cid cid.Cid) (ProtoUnmarshallable, error) {
	content, err := p.contentStorage.Read(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("read message from content storage: %w", err)
	}
	return unmarshallable{
		protoBuf: content,
		options:  proto.UnmarshalOptions{
			DiscardUnknown:    true,
		},
	}, nil
}

func (p *ProtoMessageStorage) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	encoded, err := p.marshalOpts.Marshal(message)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("marshall message: %w", err)
	}

	hash, err := multihash.Sum(encoded, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("get SHA256 multihash sum from mashalled message: %s", err)
	}
	cidValue := cid.NewCidV1(multihash.SHA2_256, hash)
	if err := p.contentStorage.Write(ctx, encoded, cidValue); err != nil {
		return cid.Cid{}, fmt.Errorf("write message to content storage: %w", err)
	}
	return cidValue, nil
}

type unmarshallable struct{
	protoBuf []byte
	options proto.UnmarshalOptions
}

func (u unmarshallable) Unmarshall(message proto.Message) error {
	return u.options.Unmarshal(u.protoBuf, message)
}

