package validator

import (
	"context"

	oppb "github.com/areknoster/public-distributed-commit-log/cmd/sentinel/open-pollution/openpollutionpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/ipfs/go-cid"
)

type SchemaValidator struct {
	messageStorage storage.MessageStorage
}

func NewSchemaValidator(messageStorage storage.MessageStorage) *SchemaValidator {
	return &SchemaValidator{messageStorage: messageStorage}
}

func (s *SchemaValidator) Validate(ctx context.Context, cid cid.Cid) error {
	message := &oppb.Message{}
	if err := s.messageStorage.Read(ctx, cid, message); err != nil {
		return sentinel.ErrorValidation{
			Kind: sentinel.ErrorValidationKindNotFound,
			Err:  err,
		}
	}
	return nil
}
