package test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/storage"
	. "github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

func TestCodecs(t *testing.T) {
	codecs := map[string]storage.Codec{
		"protoBuf": ProtoBuf{},
		"json":     Json{},
	}

	for name, codec := range codecs {
		t.Run(name, func(t *testing.T) {
			message := &testpb.Message{
				IdIncremental: 1245,
				Uuid:          uuid.NewString(),
				Created:       timestamppb.Now(),
			}
			encodedMessage, err := codec.Encode(message)
			require.NoError(t, err)

			unmarshallable := codec.Decode(encodedMessage)
			gotMesssage := &testpb.Message{}
			require.NoError(t, unmarshallable.Decode(gotMesssage))
			assert.True(t, proto.Equal(message, gotMesssage))
		})
	}
}
