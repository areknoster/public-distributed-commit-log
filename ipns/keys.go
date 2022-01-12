package ipns

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	ipfscrypto "github.com/libp2p/go-libp2p-core/crypto"
)

func ReadECKeys(privKeyPath string) (ipfscrypto.PrivKey, ipfscrypto.PubKey, error) {
	pemContent, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read %s file content: %w", privKeyPath, err)
	}

	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}
	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse PEM block containing the private key: %w", err)
	}
	return ipfscrypto.KeyPairFromStdKey(priv)
}
