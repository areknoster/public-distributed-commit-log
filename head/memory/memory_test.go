package memory

import (
	"github.com/areknoster/public-distributed-commit-log/head"
)

var _ head.Manager = &HeadManager{}
