package ipns

import (
	"fmt"
	"path"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipns"
	ipfscrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/rs/zerolog/log"
)

type Manager interface {
	UpdateIPNSEntry(string) error
	GetIPNSAddr() string
}

type NopManager struct{}

func NewNopManager() *NopManager {
	return &NopManager{}
}

func (m *NopManager) UpdateIPNSEntry(commitCID string) error {
	log.Debug().Msg("updating ipns entry in nop manager")
	return nil
}

func (m *NopManager) GetIPNSAddr() string {
	return ""
}

type IPNSManager struct {
	privKey ipfscrypto.PrivKey
	pubKey  ipfscrypto.PubKey
	shell   *shell.Shell

	ipnsAddr string
}

func NewIPNSManager(privKey ipfscrypto.PrivKey,
	pubKey ipfscrypto.PubKey,
	shell *shell.Shell) *IPNSManager {
	return &IPNSManager{
		privKey: privKey,
		pubKey:  pubKey,
		shell:   shell,
	}
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
