package messagestorage

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/pdcl"
	"github.com/areknoster/public-distributed-commit-log/storage"
)

// contentStorageWrapper creates MessageStorage based on ContentStorage implementation
type contentStorageWrapper struct {
	storage.MessageReader
	storage.MessageWriter
}

func NewContentStorageWrapper(contentStorage storage.ContentStorage, codec storage.Codec) storage.MessageStorage {
	return &contentStorageWrapper{
		MessageReader: NewContentReaderWrapper(contentStorage, codec),
		MessageWriter: NewContentWriterWrapper(contentStorage, codec),
	}
}

type contentReaderWrapper struct {
	reader  storage.ContentReader
	decoder storage.Decoder
}

func NewContentReaderWrapper(contentReader storage.ContentReader, decoder storage.Decoder) storage.MessageReader {
	return &contentReaderWrapper{
		reader:  contentReader,
		decoder: decoder,
	}
}

func (p *contentReaderWrapper) Read(ctx context.Context, cid cid.Cid) (storage.ProtoDecodable, error) {
	content, err := p.reader.Read(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("read message from content storage: %w", err)
	}
	return p.decoder.Decode(content), nil
}

type contentWriterWrapper struct {
	writer  storage.ContentWriter
	encoder storage.Encoder
}

func NewContentWriterWrapper(writer storage.ContentWriter, encoder storage.Encoder) *contentWriterWrapper {
	return &contentWriterWrapper{writer: writer, encoder: encoder}
}

func (p *contentWriterWrapper) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	encoded, err := p.encoder.Encode(message)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("marshall message: %w", err)
	}

	messageCID, err := pdcl.CID(encoded)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("get CID from mashalled message: %s", err)
	}
	if err := p.writer.Write(ctx, encoded, messageCID); err != nil {
		return cid.Cid{}, fmt.Errorf("write message to content storage: %w", err)
	}
	return messageCID, nil
}
