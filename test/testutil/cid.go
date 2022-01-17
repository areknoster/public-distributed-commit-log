// Package testutils defines helpers commonly used in testing.
package testutil

import (
	"math/rand"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"

	"github.com/areknoster/public-distributed-commit-log/pdcl"
)

func RandomCID(t *testing.T) cid.Cid {
	randBytes := make([]byte, 128)
	_, err := rand.Read(randBytes)
	require.NoError(t, err)
	c, err := pdcl.CID(randBytes)
	require.NoError(t, err)
	return c
}
