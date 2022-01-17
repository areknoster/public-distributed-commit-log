// Package ipns enables interaction with InterPlanetary Name System.
package ipns

import (
	"crypto"
	"crypto/ed25519"
	"fmt"
	"path"
	"sync"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipns"
	ipfscrypto "github.com/libp2p/go-libp2p-core/crypto"
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
	privKey ipfscrypto.PrivKey
	pubKey  ipfscrypto.PubKey
	shell   *shell.Shell

	ipnsAddr string
}

func NewIPNSManager(privKey crypto.PrivateKey, shell *shell.Shell) (*IPNSManager, error) {
	// todo: remove this after https://github.com/libp2p/go-libp2p-core/pull/234 is merged
	if val, isEd := privKey.(ed25519.PrivateKey); isEd {
		privKey = &val
	}

	priv, pub, err := ipfscrypto.KeyPairFromStdKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("get ipfscrypto key pair from private key: %w", err)
	}

	return &IPNSManager{
		privKey: priv,
		pubKey:  pub,
		shell:   shell,
	}, nil
}

func (m *IPNSManager) UpdateIPNSEntry(commitCID string) error {
	ipfsAddr := path.Join("/ipfs/", commitCID)
	ipnsRecord, err := ipns.Create(m.privKey, []byte(ipfsAddr), 0, time.Time{}, 24*time.Hour)
	if err != nil {
		return err
	}
	ipns.EmbedPublicKey(m.pubKey, ipnsRecord)

	fmt.Println("publishing...")
	resp, err := m.shell.PublishWithDetails(ipfsAddr, "", time.Hour, time.Hour, false)
	if err != nil {
		return fmt.Errorf("publishing ipns update to ipfs daemon: %v", err)
	}
	fmt.Printf("SUCCESS: %+v", resp)
	m.ipnsAddr = resp.Name
	return nil
}

func (m *IPNSManager) GetIPNSAddr() string {
	return m.ipnsAddr
}
