package ethereum

// SignatureLength is the expected length of the Signature.
const SignatureLength = 65

// Signature represents the 65 byte signature.
type Signature [SignatureLength]byte

func SignatureFromBytes(b []byte) Signature {
	var s Signature
	copy(s[:], b)
	return s
}

func SignatureFromVRS(v uint8, r [32]byte, s [32]byte) Signature {
	return SignatureFromBytes(append(append(append([]byte{}, r[:]...), s[:]...), v))
}

func (s Signature) VRS() (sv uint8, sr [32]byte, ss [32]byte) {
	copy(sr[:], s[:32])
	copy(ss[:], s[32:64])
	sv = s[64]
	return
}

func (s Signature) Bytes() []byte {
	return s[:]
}

type Signer interface {
	// Address returns account's address used to sign data. May be empty if
	// the signer is used only to verify signatures.
	Address() Address
	// SignTransaction signs transaction. Signed transaction will be set
	// to the SignedTx field in the Transaction structure.
	SignTransaction(transaction *Transaction) error
	// Signature signs the hash of the given data and returns it.
	Signature(data []byte) (Signature, error)
	// Recover returns the wallet address that created the given signature.
	Recover(signature Signature, data []byte) (*Address, error)
}
