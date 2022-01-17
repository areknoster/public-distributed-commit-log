// Package pdcl defines common logic that applies to entire module.
package pdcl

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

// ParseCID parses cid from string but handles gracefully empty cid string as cid.Undef
func ParseCID(cidStr string) (cid.Cid, error) {
	if cidStr == "" || cidStr == cid.Undef.String() {
		return cid.Undef, nil
	}
	return cid.Decode(cidStr)
}

func CID(bin []byte) (cid.Cid, error) {
	hash, err := multihash.Sum(bin, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("get SHA256 multihash sum from mashalled message: %s", err)
	}
	return cid.NewCidV1(multihash.SHA2_256, hash), nil
}
