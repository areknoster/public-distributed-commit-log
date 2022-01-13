package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/storage/content/memory"
	messagestorage "github.com/areknoster/public-distributed-commit-log/storage/message"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

func TestStorageWrapper(t *testing.T) {
	contentStorage := memory.NewStorage()
	messageStorage := messagestorage.NewContentStorageWrapper(contentStorage, pbcodec.Json{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	message := &testpb.Message{
		IdIncremental: 123,
		Uuid:          uuid.NewString(),
		Created:       timestamppb.Now(),
	}
	cid, err := messageStorage.Write(ctx, message)
	require.NoError(t, err)
	gotDecodable, err := messageStorage.Read(ctx, cid)
	require.NoError(t, err)

	gotMessage := &testpb.Message{}
	require.NoError(t, gotDecodable.Decode(gotMessage))

	assert.True(t, proto.Equal(message, gotMessage))
}
