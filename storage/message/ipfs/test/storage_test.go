package test

import (
	"context"
	"testing"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/areknoster/public-distributed-commit-log/storage"
	. "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

var _ storage.MessageStorage = &DaemonStorage{}

func TestStorage_ReadWrite(t *testing.T) {
	sh := shell.NewShell("localhost:5001")
	storage := NewStorage(sh, pbcodec.Json{})
	ctx := context.Background()
	const messageID = 98327433242
	messageCID, err := storage.Write(ctx, &testpb.Message{IdIncremental: messageID})
	t.Log("cid", messageCID.String())
	require.NoError(t, err)
	unmarshallable, err := storage.Read(ctx, messageCID)
	require.NoError(t, err)
	gotMessage := &testpb.Message{}
	require.NoError(t, unmarshallable.Decode(gotMessage))
	assert.EqualValues(t, messageID, gotMessage.IdIncremental)
}
