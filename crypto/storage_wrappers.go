package pdclcrypto

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

	"github.com/areknoster/public-distributed-commit-log/pdclpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
)

// SignedMessageWriter is a decorator on top of MessageWriter that puts every message into a signed envelope.
type SignedMessageWriter struct {
	signerID string
	signer   crypto.Signer
	base     storage.MessageWriter
	codec    storage.Codec
}

func NewSignedMessageWriter(base storage.MessageWriter, codec storage.Codec, signerID string, signer crypto.Signer) *SignedMessageWriter {
	return &SignedMessageWriter{
		signerID: signerID,
		signer:   signer,
		base:     base,
		codec:    codec,
	}
}

func (s *SignedMessageWriter) Write(ctx context.Context, message proto.Message) (cid.Cid, error) {
	rawMessage, err := s.codec.Encode(message)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("encode message: %w", err)
	}

	signature, err := s.signer.Sign(rand.Reader, rawMessage, crypto.Hash(0))
	if err != nil {
		return cid.Cid{}, fmt.Errorf("create message signature: %w", err)
	}

	signedEnvelope := &pdclpb.SignedEnvelope{
		Message:   rawMessage,
		Signature: signature,
		SignerId:  s.signerID,
	}
	return s.base.Write(ctx, signedEnvelope)
}

func NewSignedMessageUnwrapper(base storage.MessageReader, decoder storage.Decoder) *SignedMessageUnwrapper {
	return &SignedMessageUnwrapper{base: base, decoder: decoder}
}

// SignedMessageUnwrapper is a decorator on top of MessageReader that unwraps message from signed envelope.
// It does not verify signature.
type SignedMessageUnwrapper struct {
	base    storage.MessageReader
	decoder storage.Decoder
}

func (s *SignedMessageUnwrapper) Read(ctx context.Context, cid cid.Cid) (storage.ProtoDecodable, error) {
	unmarshallable, err := s.base.Read(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("read from base reader: %w", err)
	}
	envelope := &pdclpb.SignedEnvelope{}
	if err := unmarshallable.Decode(envelope); err != nil {
		return nil, fmt.Errorf("unmarshal to signed envelope: %w", err)
	}

	return s.decoder.Decode(envelope.Message), nil
}

type SignedMessageVerifyReader struct {
	base     storage.MessageReader
	registry SignerRegistry
	decoder  storage.Decoder
}

func NewSignedMessageVerifyReader(base storage.MessageReader, decoder storage.Decoder, registry SignerRegistry) *SignedMessageVerifyReader {
	return &SignedMessageVerifyReader{base: base, registry: registry, decoder: decoder}
}

var (
	// VerificationErr indicates that there was an error during message verification.
	//  It is deliberately vague to avoid adaptive attacks.
	VerificationErr = errors.New("message not verified successfully")

	SignerNotInRegistry = errors.New("signer not found in registry")
)

func (s *SignedMessageVerifyReader) Read(ctx context.Context, cid cid.Cid) (storage.ProtoDecodable, error) {
	unmarshallable, err := s.base.Read(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("read from base reader: %w", err)
	}
	envelope := new(pdclpb.SignedEnvelope)
	if err := unmarshallable.Decode(envelope); err != nil {
		return nil, fmt.Errorf("unmarshal to signed envelope: %w", err)
	}

	if err := s.verifyEnvelope(envelope); err != nil {
		return nil, err
	}

	return s.decoder.Decode(envelope.Message), nil
}

func (s *SignedMessageVerifyReader) verifyEnvelope(envelope *pdclpb.SignedEnvelope) error {
	pubKey := s.registry.Get(envelope.SignerId)

	switch pk := pubKey.(type) {
	case nil:
		return SignerNotInRegistry
	case ed25519.PublicKey:
		return s.verifiedBool(ed25519.Verify(pk, envelope.Message, envelope.Signature))
	case *ecdsa.PublicKey:
		return s.verifiedBool(ecdsa.VerifyASN1(pk, envelope.Message, envelope.Signature))
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
