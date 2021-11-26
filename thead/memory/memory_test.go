package memory

import (
	"github.com/areknoster/public-distributed-commit-log/thead"
)

var _ thead.Manager = &HeadManager{}
