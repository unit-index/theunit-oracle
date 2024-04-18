package ethkey

import (
	"bytes"
	"errors"

	"github.com/libp2p/go-libp2p-core/crypto"
	cryptoPB "github.com/libp2p/go-libp2p-core/crypto/pb"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

type PubKey struct {
	address ethereum.Address
}

func NewPubKey(address ethereum.Address) crypto.PubKey {
	return &PubKey{
		address: address,
	}
}

// Bytes implements the crypto.Key interface.
func (p *PubKey) Bytes() ([]byte, error) {
	return crypto.MarshalPublicKey(p)
}

// Equals implements the crypto.Key interface.
func (p *PubKey) Equals(key crypto.Key) bool {
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
func (p *PubKey) Raw() ([]byte, error) {
	return p.address[:], nil
}

// Type implements the crypto.Key interface.
func (p *PubKey) Type() cryptoPB.KeyType {
	return KeyTypeID
}

// Verify implements the crypto.PubKey interface.
func (p *PubKey) Verify(data []byte, sig []byte) (bool, error) {
	// Fetch public address from signature:
	addr, err := NewSigner().Recover(ethereum.SignatureFromBytes(sig), data)
	if err != nil {
		return false, err
	}

	// Verify address:
	return bytes.Equal(addr.Bytes(), p.address[:]), nil
}

// UnmarshalEthPublicKey returns a public key from input bytes.
func UnmarshalEthPublicKey(data []byte) (crypto.PubKey, error) {
	if len(data) != ethereum.AddressLength {
		return nil, errors.New("expect eth public key data size to be 20")
	}

	var addr ethereum.Address
	copy(addr[:], data)
	return &PubKey{address: addr}, nil
}
