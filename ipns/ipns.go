package ipns

import (
	"crypto"
	"fmt"
	"io"
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
	ResolveIPNS(string) (string, error)
}

type TestManager struct {
	resolved string
	mu       sync.RWMutex
}

func NewTestManager() *TestManager {
	return &TestManager{}
}

func (m *TestManager) UpdateIPNSEntry(commitCID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resolved = commitCID
	return nil
}

func (m *TestManager) GetIPNSAddr() string {
	return ""
}

func (m *TestManager) ResolveIPNS(_ string) (string, error) {
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

func (m *IPNSManager) CatIPNS(filename string) ([]byte, error) {
	bytes, err := m.shell.Cat(path.Join("/ipns/", filename))
	if err != nil {
		return nil, fmt.Errorf("cat %s from IPNS: %w", filename, err)
	}
	return io.ReadAll(bytes)
}

func (m *IPNSManager) ResolveIPNS(filename string) (string, error) {
	cid, err := m.shell.Resolve(path.Join("/ipns/", filename))
	if err != nil {
		return "", fmt.Errorf("resolve %s from IPNS: %w", filename, err)
	}
	return cid, nil
}
