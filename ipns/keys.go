package ipns

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	ipfscrypto "github.com/libp2p/go-libp2p-core/crypto"
)

func ParseKeyPair(pemContent []byte) (ipfscrypto.PrivKey, ipfscrypto.PubKey, error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse PEM block containing the private key: %w", err)
	}

	// todo: remove this after it's merged/resolved: https://github.com/libp2p/go-libp2p-core/pull/234
	if val, isEd := priv.(ed25519.PrivateKey); isEd {
		priv = &val
	}

	return ipfscrypto.KeyPairFromStdKey(priv)
}

func ReadKeyPair(privKeyPath string) (ipfscrypto.PrivKey, ipfscrypto.PubKey, error) {
	pemContent, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read %s file content: %w", privKeyPath, err)
	}

	return ParseKeyPair(pemContent)
}
