package validator

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	pdclcrypto "github.com/areknoster/public-distributed-commit-log/crypto"
	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

type Config struct {
	// Retries sets the number of retries done when reading message from storage.
	// They are run with exponential backoff, with 1 second to and 2 exponent,
	// so it is reasonable for it to stay below 10
	Retries uint `envconfig:"RETRIES" default:"6"`
}

type SignedValidator struct {
	reader storage.MessageReader
	config Config
}

func New(messageReader storage.MessageReader, decoder storage.Decoder, config Config) (*SignedValidator, error) {
	registry, err := loadRegistry()
	if err != nil {
		return nil, fmt.Errorf("load producer registry: %w", err)
	}

	signedMessageReader := pdclcrypto.NewSignedMessageVerifyReader(messageReader, decoder, registry)

	return &SignedValidator{
		reader: signedMessageReader,
		config: config,
	}, nil
}

func loadRegistry() (pdclcrypto.MemorySignerRegistry, error) {
	registry := pdclcrypto.MemorySignerRegistry{}
	for _, p := range producers {
		block, _ := pem.Decode([]byte(p.publicKey))
		if block == nil {
			return nil, fmt.Errorf("producer_id=%s, failed to parse PEM block containing the public key. ", p.id)
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("producer_id=%s, failed to parse PEM block containing the public key: %w", p.id, err)
		}
		registry[p.id] = pub
	}
	return registry, nil
}

func (s *SignedValidator) Validate(ctx context.Context, cid cid.Cid) error {
	unmarshallable, err := s.readMessage(ctx, cid)
	if err != nil {
		return err
	}

	msg := new(testpb.Message)
	if err := unmarshallable.Decode(msg); err != nil {
		return sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindIncorrectContent,
			Err:  fmt.Errorf("can't deserialize with test message proto schema: %w", err),
		}
	}

	if err := s.validateCreatedDate(msg.Created); err != nil {
		return err
	}

	if _, err := uuid.Parse(msg.Uuid); err != nil {
		return sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindIncorrectContent,
			Err:  fmt.Errorf("incorrect UUID field: %w", err),
		}
	}
	return nil
}

func (s *SignedValidator) readMessage(ctx context.Context, cid cid.Cid) (storage.ProtoDecodable, error) {
	var (
		err            error
		unmarshallable storage.ProtoDecodable
		logger         = log.With().Stringer("cid", cid).Logger()
	)
	for i := uint(0); i <= s.config.Retries; i++ {
		unmarshallable, err = s.reader.Read(ctx, cid)
		if err == nil {
			break
		}
		logger.Info().Err(err).Msg("read message")
		time.Sleep(time.Second * (1 << i))
	}

	switch {
	case err == nil:
		// continue peacefully
	case errors.Is(err, storage.ErrInternal), errors.Is(err, storage.ErrTimeout):
		return nil, sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindInternal,
			Err:  err,
		}
	case errors.Is(err, storage.ErrNotFound):
		return nil, sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindNotFound,
			Err:  err,
		}
	default:
		return nil, sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindUnknown,
			Err:  err,
		}
	}
	return unmarshallable, nil
}

func (s *SignedValidator) validateCreatedDate(timestamp *timestamppb.Timestamp) error {
	if timestamp == nil || !timestamp.IsValid() {
		return sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindIncorrectContent,
			Err:  errors.New("created timestamp can't be nil"),
		}
	}
	t := timestamp.AsTime()
	if t.Before(time.Now().Add(-12*time.Hour)) || t.After(time.Now()) {
		return sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindIncorrectContent,
			Err:  errors.New("message must be created within last 12 hours"),
		}
	}
	return nil
}
