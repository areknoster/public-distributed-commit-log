package pbcodec

import (
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/storage"
)

type ProtoBuf struct{}

func (p ProtoBuf) Encode(message proto.Message) (storage.RawMessage, error) {
	return proto.MarshalOptions{
		Deterministic: true,
	}.Marshal(message)
}

func (p ProtoBuf) Decode(message storage.RawMessage) storage.ProtoDecodable {
	return protobufUnmarshallable(message)
}

type protobufUnmarshallable storage.RawMessage

func (pu protobufUnmarshallable) Decode(message proto.Message) error {
	return proto.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(pu, message)
}
