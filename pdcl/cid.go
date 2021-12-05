// Package pdcl defines common logic that applies to entire module
package pdcl

import "github.com/ipfs/go-cid"

// ParseCID parses cid from string but handles gracefully empty cid string as cid.Undef
func ParseCID(cidStr string) (cid.Cid, error) {
	if cidStr == "" || cidStr == cid.Undef.String() {
		return cid.Undef, nil
	}
	return cid.Decode(cidStr)
}
