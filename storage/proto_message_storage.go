package storage

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/pdcl"
)

// ProtoMessageStorage creates MessageStorage based on ContentStorage implementation
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

func GetMarshallOpts() proto.MarshalOptions {
	return marshallOpts
}

var marshallOpts = proto.MarshalOptions{
	Deterministic: true,
}

func ProtoEncode(message proto.Message) ([]byte, error) {
	return marshallOpts.Marshal(message)
}

func ProtoDecode(content []byte) ProtoUnmarshallable {
	return protobufUnmarshallable{
		protoBuf: content,
		options: proto.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
}

type protobufUnmarshallable struct {
	protoBuf []byte
	options  proto.UnmarshalOptions
}

func (u protobufUnmarshallable) Unmarshall(message proto.Message) error {
	return u.options.Unmarshal(u.protoBuf, message)
}
