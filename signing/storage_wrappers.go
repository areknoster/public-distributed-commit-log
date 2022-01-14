package signing

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/ipfs/go-cid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/areknoster/public-distributed-commit-log/pdclpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
)

// SignedMessageWriter is a decorator on top of MessageWriter that puts every message into a signed envelope.
type SignedMessageWriter struct {
	signerID string
	signer   crypto.Signer
	base     storage.MessageWriter
}

func NewSignedMessageWriter(base storage.MessageWriter, signerID string, signer crypto.Signer) *SignedMessageWriter {
	return &SignedMessageWriter{signerID: signerID, signer: signer, base: base}
}

func (s *SignedMessageWriter) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	any := new(anypb.Any)
	if err := anypb.MarshalFrom(any, message, storage.GetMarshallOpts()); err != nil {
		return cid.Cid{}, fmt.Errorf("marshal to anypb: %w", err)
	}

	signature, err := s.signer.Sign(rand.Reader, any.Value, crypto.Hash(0))
	if err != nil {
		return cid.Cid{}, fmt.Errorf("create message signature: %w", err)
	}

	signedEnvelope := &pdclpb.SignedEnvelope{
		Message:   any,
		Signature: signature,
		SignerId:  s.signerID,
	}
	return s.base.Write(ctx, signedEnvelope)
}

// SignedMessageUnwrapper is a decorator on top of MessageReader that unwraps message from signed envelope.
// It does not verify signature.
type SignedMessageUnwrapper struct {
	Base storage.MessageReader
}

func (s *SignedMessageUnwrapper) Read(ctx context.Context, cid cid.Cid) (storage.ProtoUnmarshallable, error) {
	unmarshallable, err := s.Base.Read(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("read from base reader: %w", err)
	}
	envelope := new(pdclpb.SignedEnvelope)
	if err := unmarshallable.Unmarshall(envelope); err != nil {
		return nil, fmt.Errorf("unmarshal to signed envelope: %w", err)
	}

	return storage.ProtoDecode(envelope.Message.Value), nil
}

type SignedMessageVerifyReader struct {
	base     storage.MessageReader
	registry SignerRegistry
}

func NewSignedMessageVerifyReader(base storage.MessageReader, registry SignerRegistry) *SignedMessageVerifyReader {
	return &SignedMessageVerifyReader{base: base, registry: registry}
}

var (
	// VerificationErr indicates that there was an error during message verification.
	//  It is deliberately vague to avoid adaptive attacks.
	VerificationErr = errors.New("message not verified successfully")

	SignerNotInRegistry = errors.New("signer not found in registry")
)

func (s *SignedMessageVerifyReader) Read(ctx context.Context, cid cid.Cid) (storage.ProtoUnmarshallable, error) {
	unmarshallable, err := s.base.Read(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("read from base reader: %w", err)
	}
	envelope := new(pdclpb.SignedEnvelope)
	if err := unmarshallable.Unmarshall(envelope); err != nil {
		return nil, fmt.Errorf("unmarshal to signed envelope: %w", err)
	}

	if err := s.verifyEnvelope(envelope); err != nil {
		return nil, err
	}

	return storage.ProtoDecode(envelope.Message.Value), nil
}

func (s *SignedMessageVerifyReader) verifyEnvelope(envelope *pdclpb.SignedEnvelope) error {
	pubKey := s.registry.Get(envelope.SignerId)
	switch pk := pubKey.(type) {
	case nil:
		return SignerNotInRegistry
	case ed25519.PublicKey:
		return s.verifiedBool(ed25519.Verify(pk, envelope.Message.Value, envelope.Signature))
	case *ecdsa.PublicKey:
		return s.verifiedBool(ecdsa.VerifyASN1(pk, envelope.Message.Value, envelope.Signature))
	default:
		return fmt.Errorf("unsupported singer key type: %T", pk)
	}
}

func (s SignedMessageVerifyReader) verifiedBool(ok bool) error {
	if !ok {
		return VerificationErr
	}
	return nil
}
