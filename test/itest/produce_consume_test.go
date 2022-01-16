package itest

import (
	"context"
	"testing"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/stretchr/testify/suite"

	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

type MemoryProduceConsumeTestSuite struct {
	ProduceConsumeTestSuite
}

func TestMemoryProduceConsumeTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryProduceConsumeTestSuite))
}

type LocalDaemonProduceConsumeTestSuite struct {
	ProduceConsumeTestSuite
}

func (s *LocalDaemonProduceConsumeTestSuite) SetupSuite() {
	s.setupMessageStorage()
	s.ProduceConsumeTestSuite.SetupSuite()
	s.waitForDaemon()
}

func (s *LocalDaemonProduceConsumeTestSuite) waitForDaemon() {
	s.T().Log("wait for daemon to start responding")
	for i := 0; i < 5; i++ {
		if _, err := s.messageStorage.Write(context.Background(), &testpb.Message{}); err == nil {
			return
		}
		time.Sleep(3 * time.Second)
	}
	s.Fail("daemon doesn't respond")
}

func (s *LocalDaemonProduceConsumeTestSuite) setupMessageStorage() {
	sh := shell.NewShell("localhost:5001")
	s.messageStorage = ipfsstorage.NewStorage(sh, pbcodec.Json{})
}

func TestLocalDaemonProduceConsumeTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ")
	}
	suite.Run(t, new(LocalDaemonProduceConsumeTestSuite))
}
