package storage

import (
	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"
)

type Encoder interface {
	Encode(message proto.Message) ([]byte, cid.Cid, error)
}

type Decoder interface{
	Decode(content []byte) (ProtoUnmarshallable, error)
}


