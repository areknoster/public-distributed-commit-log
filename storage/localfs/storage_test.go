package localfs

import "github.com/areknoster/public-distributed-commit-log/pkg/storage"

var _ storage.Storage = &Storage{}
