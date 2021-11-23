// Package cid is a wrapper over github.com/ipfs/go-cid, with defaults and additional empty cid handling
package pdcl

import  "github.com/ipfs/go-cid"

func ParseCID(cidStr string )  (cid.Cid, error){
	if cidStr == "" || cidStr == cid.Undef.String(){
		return cid.Undef, nil
	}
	return cid.Decode(cidStr)
}