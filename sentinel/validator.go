package sentinel

import (
	"context"

	"github.com/ipfs/go-cid"
)

type Validator interface {
	Validate(ctx context.Context, cid cid.Cid) error
}

type ErrorValidationKind uint

const (
	ErrorValidationKindUnknown ErrorValidationKind = iota
	ErrorValidationKindNotFound
	ErrorValidationKindInternal
)

type ErrorValidation struct {
	Kind ErrorValidationKind
	Err  error
}

func (e ErrorValidation) Error() string {
	return e.Err.Error()
}

func (e ErrorValidation) Unwrap() error {
	return e.Err
}
