package ipns

import (
	"fmt"
	"path"
	"strings"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"

	"github.com/areknoster/public-distributed-commit-log/pdcl"
)

type resolver interface {
	// resolveIPNS finds IPFS address that's pointed by given IPNS address
	resolveIPNS(ipnsName string) (ipfsAddress cid.Cid, err error)
}

type shellResolver struct {
	shell *shell.Shell
}

func newShellResolver(shell *shell.Shell) *shellResolver {
	return &shellResolver{shell: shell}
}

func (m *shellResolver) resolveIPNS(ipnsName string) (ipfsAddress cid.Cid, err error) {
	resolvedAddr, err := m.shell.Resolve(path.Join("/ipns/", ipnsName))
	if err != nil {
		return cid.Undef, fmt.Errorf("resolve %s from IPNS: %w", ipnsName, err)
	}
	resolvedCid, err := pdcl.ParseCID(strings.TrimPrefix(resolvedAddr, "/ipfs/"))
	if err != nil {
		return cid.Undef, fmt.Errorf("parse resolved IPNS address to CID: %w", err)
	}
	return resolvedCid, nil
}
