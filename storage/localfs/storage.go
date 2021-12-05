package localfs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	dirPath string
}

func (s *Storage) Read(ctx context.Context, cid cid.Cid) ([]byte, error) {
	filePath := path.Join(s.dirPath, cid.String())
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file with message: %w", err)
	}
	return content, nil
}

func (s *Storage) Write(ctx context.Context, content []byte, cidValue cid.Cid) error {
	filePath := path.Join(s.dirPath, cidValue.String())
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		if fileInfo.Size() != int64(len(content)) {
			return fmt.Errorf("content with given CID is already saved, but it has different length")
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("unknown file stat error: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	if _, err := io.Copy(file, bytes.NewBuffer(content)); err != nil {
		if err := file.Close(); err != nil {
			log.Error().Str("file", filePath).Err(err).Msg("close file")
		}
		return fmt.Errorf("copy message content to file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close file: %w", err)
	}
	return nil
}

func NewStorage(dirPath string) (*Storage, error) {
	err := createDirIfNotExists(dirPath)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		dirPath: dirPath,
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
