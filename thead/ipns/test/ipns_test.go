// Package ipns provides implementation for pdcl head management based on ipns protocol.
package ipns

import (
	"context"
	"testing"
	"time"

	"github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
	"github.com/areknoster/public-distributed-commit-log/thead"
	. "github.com/areknoster/public-distributed-commit-log/thead/ipns"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadManager_WriteReadHead(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping because ipfs daemon is needed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	sh := shell.NewShell("localhost:5001")

	storage := ipfs.NewStorage(sh, pbcodec.ProtoBuf{})
	msg := testpb.MakeCurrentRandomTestMessage()
	id, err := storage.Write(ctx, msg)
	require.NoError(t, err)
	t.Log("message written")

	var headManager thead.Manager
	headManager, err = NewBasicHeadManager(sh)
	require.NoError(t, err)
	require.NoError(t, headManager.SetHead(ctx, id))

	headCid, err := headManager.ReadHead(ctx)
	require.NoError(t, err)
	assert.Equal(t, id, headCid)
}
