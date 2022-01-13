package pdclcrypto

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func ReadEd25519(privKeyPath string) (crypto.Signer, error) {
	pemContent, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read %s file content: %w", privKeyPath, err)
	}

	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key: %w", err)
	}
	key, ok := priv.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not ed25519.PrivateKey but %T", priv)
	}
	return key, nil
}
