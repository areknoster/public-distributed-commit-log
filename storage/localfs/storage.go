package localfs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Storage struct {
	dirPath     string
	marshalOpts proto.MarshalOptions
}

func NewStorage(dirPath string) (*Storage, error) {
	err := createDirIfNotExists(dirPath)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		dirPath: dirPath,
		marshalOpts: proto.MarshalOptions{
			Deterministic: true,
		},
	}
	return s, nil
}

func createDirIfNotExists(dirPath string) error {
	fileInfo, err := os.Stat(dirPath)
	if err == nil {
		if !fileInfo.IsDir() {
			return fmt.Errorf("path does not point to directory")
		}
		return nil // directory exists
	}

	if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dirPath, fs.ModeDir|0o755); err != nil {
		return fmt.Errorf("create directory %s: %w", dirPath, err)
	}
	return nil
}

func (s *Storage) Read(ctx context.Context, cid cid.Cid, message proto.Message) error {
	filePath := path.Join(s.dirPath, cid.String())
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file with message: %w", err)
	}
	if err := proto.Unmarshal(content, message); err != nil {
		return fmt.Errorf("%w: unmarshal found message to given structure: %s", storage.ErrUnmarshall, err)
	}
	return nil
}

func (s *Storage) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	encoded, err := s.marshalOpts.Marshal(message)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("%w: %s", storage.ErrInternal, err.Error())
	}

	hash, err := multihash.Sum(encoded, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("%w: get SHA256 multihash sum from mashalled message: %s", storage.ErrInternal, err.Error())
	}
	cidValue := cid.NewCidV1(multihash.SHA2_256, hash)
	cidValue.String()

	filePath := path.Join(s.dirPath, cidValue.String())
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		if fileInfo.Size() != int64(len(encoded)) {
			return cid.Cid{}, fmt.Errorf("%w: content with given CID is already saved, but it has different length", storage.ErrInternal)
		}
		return cidValue, nil
	}
	if !os.IsNotExist(err) {
		return cid.Cid{}, fmt.Errorf("%w: unknown file stat error: %s", storage.ErrInternal, err.Error())
	}

	file, err := os.Create(filePath)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("create file: %w", err)
	}
	if _, err := io.Copy(file, bytes.NewBuffer(encoded)); err != nil {
		if err := file.Close(); err != nil {
			log.Error().Str("file", filePath).Err(err).Msg("close file")
		}
		return cid.Cid{}, fmt.Errorf("copy message content to file: %w", err)
	}
	if err := file.Close(); err != nil {
		return cid.Cid{}, fmt.Errorf("close file: %w", err)
	}
	return cidValue, nil
}
