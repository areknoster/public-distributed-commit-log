package ipfs

import "github.com/areknoster/public-distributed-commit-log/storage"

var _ storage.MessageStorage = &Storage{}
