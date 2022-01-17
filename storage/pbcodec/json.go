// Package pbcoded provides codecs for (de)serializing messages.
package pbcodec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/storage"
)

type Json struct{}

func (j Json) Encode(message proto.Message) (storage.RawMessage, error) {
	return protojson.MarshalOptions{
		UseEnumNumbers: true, // this preserves backward compatibility
	}.Marshal(message)
}

func (j Json) Decode(message storage.RawMessage) storage.ProtoDecodable {
	return jsonUnmarshallable(message)
}

type jsonUnmarshallable storage.RawMessage

func (ju jsonUnmarshallable) Decode(message proto.Message) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(ju, message)
}
