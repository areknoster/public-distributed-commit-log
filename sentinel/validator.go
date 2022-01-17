// Package sentinel defines type interfaces and other code used by sentinel.
package sentinel

import (
	"context"

	"github.com/ipfs/go-cid"
)

// Validator validates message contents in sentinel. It should be implemented by topic owner.
type Validator interface {
	Validate(ctx context.Context, cid cid.Cid) error
}

type ErrorValidationKind uint

const (
	ErrorValidationKindUnknown ErrorValidationKind = iota
	ErrorValidationKindNotFound
	ErrorValidationKindInternal
	ErrorValidationKindIncorrectContent
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
