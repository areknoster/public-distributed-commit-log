package test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/areknoster/public-distributed-commit-log/internal/testpb"
	. "github.com/areknoster/public-distributed-commit-log/signing"
	"github.com/areknoster/public-distributed-commit-log/storage"
	memorystorage "github.com/areknoster/public-distributed-commit-log/storage/memory"
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
	store := storage.NewProtoMessageStorage(memorystorage.New())
	message := &testpb.Message{
		IdIncremental: 1,
	}

	t.Run("if signer is not found in singers registry, SignerNotInRegistry should be wrapped", func(t *testing.T) {
		signingWriter := NewSignedMessageWriter(store, "not-known-signer-id", privKey1)
		verifyingReader := NewSignedMessageVerifyReader(store, registry)
		ccid, err := signingWriter.Write(ctx, message)
		require.NoError(t, err)
		_, err = verifyingReader.Read(ctx, ccid)
		assert.ErrorIs(t, err, SignerNotInRegistry)
	})

	t.Run("if signer is not found in singers registry, SignerNotInRegistry should be wrapped", func(t *testing.T) {
		signingWriter := NewSignedMessageWriter(store, "not-known-signer-id", privKey1)
		verifyingReader := NewSignedMessageVerifyReader(store, registry)
		ccid, err := signingWriter.Write(ctx, message)
		require.NoError(t, err)
		_, err = verifyingReader.Read(ctx, ccid)
		assert.ErrorIs(t, err, SignerNotInRegistry)
	})

	t.Run("if different key was used for signing than is registered, then verification error should be wrapped", func(t *testing.T) {
		signingWriter := NewSignedMessageWriter(store, signerID, privKey2)
		verifyingReader := NewSignedMessageVerifyReader(store, registry)
		ccid, err := signingWriter.Write(ctx, message)
		require.NoError(t, err)
		_, err = verifyingReader.Read(ctx, ccid)
		assert.ErrorIs(t, err, VerificationErr)
	})

	t.Run("happy path: message is signed, verified and unmarshallable unmarshalls to message content", func(t *testing.T) {
		signingWriter := NewSignedMessageWriter(store, signerID, privKey1)
		verifyingReader := NewSignedMessageVerifyReader(store, registry)
		ccid, err := signingWriter.Write(ctx, message)
		require.NoError(t, err)
		unmarshallable, err := verifyingReader.Read(ctx, ccid)
		require.NoError(t, err)

		got := &testpb.Message{}
		require.NoError(t, unmarshallable.Unmarshall(got))

		assert.True(t, proto.Equal(message, got), "got same message as produced")
	})
}
