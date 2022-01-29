// Package ipns enables interaction with InterPlanetary Name System.
package ipns

import (
	"fmt"
	"path"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multibase"
)

type Manager interface {
	UpdateIPNSEntry(string) error
	GetIPNSAddr() string
}

type TestManagerResolver struct {
	resolved string
	mu       sync.RWMutex
}

func NewTestManager() *TestManagerResolver {
	return &TestManagerResolver{}
}

func (m *TestManagerResolver) UpdateIPNSEntry(commitCID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resolved = commitCID
	return nil
}

func (m *TestManagerResolver) GetIPNSAddr() string {
	return ""
}

func (m *TestManagerResolver) ResolveIPNS(_ string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("/ipfs/%s", m.resolved), nil
}

type IPNSManager struct {
	shell    *shell.Shell
	ipnsAddr string
}

func getIPNSAddress(sh *shell.Shell) (string, error) {
	// this implementation is extremely non-obvious
	// because IPFS doesn't normally allow for finding IPNS address
	// of given key unless some file is added to it.
	resp, err := sh.ID()
	if err != nil {
		return "", fmt.Errorf("get IPFS ID: %w", err)
	}
	pid, err := peer.Decode(resp.ID)
	if err != nil {
		return "", fmt.Errorf("decode peer ID: %w", err)
	}
	ipnsAddr, err := peer.ToCid(pid).StringOfBase(multibase.Base36)
	if err != nil {
		return "", fmt.Errorf("encode ipns address: %w", err)
	}
	return ipnsAddr, nil
}

func NewIPNSManager(sh *shell.Shell) (*IPNSManager, error) {
	ipnsAddress, err := getIPNSAddress(sh)
	if err != nil {
		return nil, fmt.Errorf("get IPNS address from IPFS node: %w", err)
	}

	return &IPNSManager{
		shell:    sh,
		ipnsAddr: ipnsAddress,
	}, nil
}

func (m *IPNSManager) UpdateIPNSEntry(commitCID string) error {
	ipfsAddr := path.Join("/ipfs/", commitCID)
	resp, err := m.shell.PublishWithDetails(ipfsAddr, "", 24*time.Hour, 10*time.Minute, false)
	if err != nil {
		return fmt.Errorf("publishing ipns update to ipfs daemon: %v", err)
	}
	m.ipnsAddr = resp.Name
	return nil
}

func (m *IPNSManager) GetIPNSAddr() string {
	return m.ipnsAddr
}
