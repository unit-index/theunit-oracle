package ethkey

import (
	"bytes"
	"errors"

	"github.com/libp2p/go-libp2p-core/crypto"
	cryptoPB "github.com/libp2p/go-libp2p-core/crypto/pb"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

type PrivKey struct {
	signer ethereum.Signer
}

func NewPrivKey(signer ethereum.Signer) crypto.PrivKey {
	return &PrivKey{
		signer: signer,
	}
}

// Bytes implements the crypto.Key interface.
func (p *PrivKey) Bytes() ([]byte, error) {
	return crypto.MarshalPrivateKey(p)
}

// Equals implements the crypto.Key interface.
func (p *PrivKey) Equals(key crypto.Key) bool {
	if p.Type() != key.Type() {
		return false
	}

	a, err := p.Raw()
	if err != nil {
		return false
	}
	b, err := key.Raw()
	if err != nil {
		return false
	}

	return bytes.Equal(a, b)
}

// Raw implements the crypto.Key interface.
func (p *PrivKey) Raw() ([]byte, error) {
	return p.signer.Address().Bytes(), nil
}

// Type implements the crypto.Key interface.
func (p *PrivKey) Type() cryptoPB.KeyType {
	return KeyTypeID
}

// Sign implements the crypto.PrivateKey interface.
func (p *PrivKey) Sign(bytes []byte) ([]byte, error) {
	s, err := p.signer.Signature(bytes)
	if err != nil {
		return nil, err
	}
	return s.Bytes(), nil
}

// GetPublic implements the crypto.PrivateKey interface.
func (p *PrivKey) GetPublic() crypto.PubKey {
	return NewPubKey(p.signer.Address())
}

// UnmarshalEthPrivateKey should return private key from input bytes, but this
// not supported for ethereum keys.
func UnmarshalEthPrivateKey(data []byte) (crypto.PrivKey, error) {
	return nil, errors.New("eth key type does not support unmarshalling")
}
