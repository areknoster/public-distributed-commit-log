package signing

import (
	"crypto"
)

type SignerRegistry interface {
	Get(signerID string) crypto.PublicKey
}

type MemorySignerRegistry map[string]crypto.PublicKey

func (m MemorySignerRegistry) Get(signerID string) crypto.PublicKey {
	return m[signerID]
}
