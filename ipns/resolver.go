package ipns

import (
	"fmt"
	"path"

	shell "github.com/ipfs/go-ipfs-api"
)

type Resolver interface {
	// ResolveIPNS finds IPFS address that's pointed by given IPNS address
	ResolveIPNS(string) (string, error)
}

type IPNSResolver struct {
	shell *shell.Shell
}

func NewIPNSResolver(shell *shell.Shell) *IPNSResolver {
	return &IPNSResolver{shell: shell}
}

func (m *IPNSResolver) ResolveIPNS(filename string) (string, error) {
	cid, err := m.shell.Resolve(path.Join("/ipns/", filename))
	if err != nil {
		return "", fmt.Errorf("resolve %s from IPNS: %w", filename, err)
	}
	return cid, nil
}
