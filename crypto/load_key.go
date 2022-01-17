// Package pdclcrypto provides storage signing wrappers and tools for parsing cryptographic keys.
package pdclcrypto

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadFromPKCSFromPEMFile(privKeyPath string) (crypto.PrivateKey, error) {
	pemContent, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read %s file content: %w", privKeyPath, err)
	}
	return ParsePKCSKeyFromPEM(pemContent)
}

func ParsePKCSKeyFromPEM(pemContent []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key: %w", err)
	}
	return priv, nil
}
