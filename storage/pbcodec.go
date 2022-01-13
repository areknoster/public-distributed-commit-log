package storage

import "google.golang.org/protobuf/proto"

type (
	// ProtoDecodable can be used to deserialize message data to proto structure
	// todo: refactor this to just passing generic type when go 1.18 is out
	ProtoDecodable interface {
		Decode(message proto.Message) error
	}
	RawMessage []byte

	// Encoder
	Encoder interface {
		Encode(message proto.Message) (RawMessage, error)
	}
	Decoder interface {
		Decode(message RawMessage) ProtoDecodable
	}

	Codec interface {
		Encoder
		Decoder
	}
)
