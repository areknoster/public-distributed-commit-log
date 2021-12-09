package consumer

import (
	"errors"
)

type Error error

var ErrContextDone Error = errors.New("context is done")
