package consumer

import (
	"context"
	"errors"
	"fmt"
	"github.com/areknoster/public-distributed-commit-log/pdcl"
	"github.com/areknoster/public-distributed-commit-log/pdclpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/ipfs/go-cid"
	"time"
)

type commit struct {
	Cid      cid.Cid
	Previous cid.Cid
	Messages []cid.Cid
}

var (
	ErrCommitNotFound = errors.New("commit not found")
)

func newCommit(itsOwnCID cid.Cid, pbCommit *pdclpb.Commit) (commit, error) {
	messagesCids := make([]cid.Cid, len(pbCommit.MessagesCids))
	var err error
	for i, messageCID := range pbCommit.MessagesCids {
		messagesCids[i], err = cid.Decode(messageCID)
		if err != nil {
			return commit{}, fmt.Errorf("decode message cid: %w")
		}
	}

	var previousCID cid.Cid
	previousCID, err = pdcl.ParseCID(pbCommit.PreviousCommitCid)
	if err != nil {
		return commit{}, fmt.Errorf("decode previous commit cid: %w")
	}

	return commit{
		Cid:      itsOwnCID,
		Previous: previousCID,
		Messages: messagesCids,
	}, nil
}

type CommitReader interface{
	GetCommit(ctx context.Context, cid cid.Cid) (commit, error)
}

type StorageCommitReader struct{
	reader storage.MessageReader
	timeout time.Duration
}

func NewStorageCommitReader(reader storage.MessageReader, timeout time.Duration) *StorageCommitReader {
	return &StorageCommitReader{reader: reader, timeout: timeout}
}

func (cr *StorageCommitReader) GetCommit(ctx context.Context, cid cid.Cid) (commit, error){
	ctxTimeout, cancel := context.WithTimeout(ctx, cr.timeout)
	defer cancel()
	unmarshallable, err := cr.reader.Read(ctxTimeout, cid)
	if err != nil {
		return commit{}, fmt.Errorf("read commit message from storage: %w", err)
	}
	pbCommit := &pdclpb.Commit{}
	if err := unmarshallable.Unmarshall(pbCommit); err != nil{
		return commit{}, fmt.Errorf("unmarshall to commit proto: %w", err)
	}

	domainCommit, err := newCommit(cid, pbCommit)
	if err != nil {
		return commit{}, fmt.Errorf("map proto commit to consumer commit: %w", err)
	}

	return domainCommit, nil
}

type FakeCommitReader struct{
	Commits map[cid.Cid]commit
}

func (f FakeCommitReader) GetCommit(ctx context.Context, cid cid.Cid) (commit, error) {
	c, exists := f.Commits[cid]
	if !exists{
		return commit{}, ErrCommitNotFound
	}
	return c, nil
}



