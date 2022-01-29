package acceptance

import (
	"testing"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multibase"
	"github.com/stretchr/testify/require"
)

func TestLol(t *testing.T) {
	sh := shell.NewShell("localhost:5001")
	resp, err := sh.ID()
	require.NoError(t, err)
	pid, err := peer.Decode(resp.ID)
	require.NoError(t, err)

	t.Log(peer.ToCid(pid).StringOfBase(multibase.Base36))
}
