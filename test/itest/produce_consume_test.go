package itest

import (
	"context"
	"testing"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/stretchr/testify/suite"

	"github.com/areknoster/public-distributed-commit-log/ipns"
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
	sh *shell.Shell
}

func (s *LocalDaemonProduceConsumeTestSuite) SetupSuite() {
	s.sh = shell.NewShell("localhost:5001")

	s.setupMessageStorage()
	s.setupIPNS()
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

func (s *LocalDaemonProduceConsumeTestSuite) setupIPNS() {
	ipnsMgr, err := ipns.NewIPNSManager(s.sh)
	s.Require().NoError(err)
	s.ipnsMgr = ipnsMgr
}

func (s *LocalDaemonProduceConsumeTestSuite) setupMessageStorage() {
	s.messageStorage = ipfsstorage.NewStorage(s.sh, pbcodec.Json{})
}

func TestLocalDaemonProduceConsumeTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ")
	}
	suite.Run(t, new(LocalDaemonProduceConsumeTestSuite))
}
