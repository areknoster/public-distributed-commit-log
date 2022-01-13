package test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	. "github.com/areknoster/public-distributed-commit-log/crypto"
	memorystorage "github.com/areknoster/public-distributed-commit-log/storage/content/memory"
	messagestorage "github.com/areknoster/public-distributed-commit-log/storage/message"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

func TestSignVerify(t *testing.T) {
	ctx := context.TODO()
	pubKey1, privKey1, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	_, privKey2, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	const signerID = "testID"
	registry := MemorySignerRegistry{
		signerID: pubKey1,
	}
	store := messagestorage.NewContentStorageWrapper(memorystorage.NewStorage(), pbcodec.Json{})
	message := &testpb.Message{
		IdIncremental: 1,
		Uuid:          uuid.NewString(),
		Created:       timestamppb.Now(),
	}

	t.Run("if signer is not found in singers registry, SignerNotInRegistry should be wrapped", func(t *testing.T) {
		signingWriter := NewSignedMessageWriter(store, pbcodec.Json{}, "not-known-signer-id", privKey1)
		verifyingReader := NewSignedMessageVerifyReader(store, pbcodec.Json{}, registry)
		ccid, err := signingWriter.Write(ctx, message)
		require.NoError(t, err)
		_, err = verifyingReader.Read(ctx, ccid)
		assert.ErrorIs(t, err, SignerNotInRegistry)
	})

	t.Run("if different key was used for signing than is registered, then verification error should be wrapped", func(t *testing.T) {
		signingWriter := NewSignedMessageWriter(store, pbcodec.Json{}, signerID, privKey2)
		verifyingReader := NewSignedMessageVerifyReader(store, pbcodec.Json{}, registry)
		ccid, err := signingWriter.Write(ctx, message)
		require.NoError(t, err)
		_, err = verifyingReader.Read(ctx, ccid)
		assert.ErrorIs(t, err, VerificationErr)
	})

	t.Run("happy path: message is signed, verified and unmarshallable unmarshalls to message content", func(t *testing.T) {
		signingWriter := NewSignedMessageWriter(store, pbcodec.Json{}, signerID, privKey1)
		verifyingReader := NewSignedMessageVerifyReader(store, pbcodec.Json{}, registry)
		ccid, err := signingWriter.Write(ctx, message)
		require.NoError(t, err)
		unmarshallable, err := verifyingReader.Read(ctx, ccid)
		require.NoError(t, err)

		got := &testpb.Message{}
		require.NoError(t, unmarshallable.Decode(got))

		assert.True(t, proto.Equal(message, got), "got same message as produced")
	})
}
