package validator

import (
	"context"

	"github.com/ipfs/go-cid"

	"github.com/areknoster/public-distributed-commit-log/cmd/openpollution/pb"
	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/storage"
)

type SchemaValidator struct {
	messageStorage storage.MessageStorage
}

func NewSchemaValidator(messageStorage storage.MessageStorage) *SchemaValidator {
	return &SchemaValidator{messageStorage: messageStorage}
}

func (s *SchemaValidator) Validate(ctx context.Context, cid cid.Cid) error {
	unmarshallable, err := s.messageStorage.Read(ctx, cid)
	if err != nil {
		return sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindNotFound,
			Err:  err,
		}
	}
	message := &pb.Message{}
	if err := unmarshallable.Unmarshall(message); err != nil {
		return sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindIncorrectContent,
			Err:  err,
		}
	}

	return nil
}
