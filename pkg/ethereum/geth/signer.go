package geth

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

var ErrInvalidSignature = errors.New("invalid Ethereum signature (V is not 27 or 28)")

type Signer struct {
	account *Account
}

// NewSigner returns a new Signer instance. If you don't want to sign any data
// and you only want to recover public addresses, you may use nil as an argument.
func NewSigner(account *Account) *Signer {
	return &Signer{
		account: account,
	}
}

// Address implements the ethereum.Signer interface.
func (s *Signer) Address() ethereum.Address {
	if s.account == nil {
		return ethereum.Address{}
	}

	return s.account.address
}

// SignTransaction implements the ethereum.Signer interface.
func (s *Signer) SignTransaction(transaction *ethereum.Transaction) error {
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:    nil,
		Nonce:      transaction.Nonce,
		GasTipCap:  transaction.PriorityFee,
		GasFeeCap:  transaction.MaxFee,
		Gas:        transaction.GasLimit.Uint64(),
		To:         &transaction.Address,
		Value:      nil,
		Data:       transaction.Data,
		AccessList: nil,
		V:          nil,
		R:          nil,
		S:          nil,
	})
	signedTx, err := s.account.wallet.SignTxWithPassphrase(
		*s.account.account,
		s.account.passphrase,
		tx,
		transaction.ChainID,
	)
	if err != nil {
		return err
	}
	transaction.SignedTx = signedTx
	return nil
}

// Signature implements the ethereum.Signer interface.
func (s *Signer) Signature(data []byte) (ethereum.Signature, error) {
	return Signature(s.account, data)
}

// Recover implements the ethereum.Signer interface.
func (s *Signer) Recover(signature ethereum.Signature, data []byte) (*ethereum.Address, error) {
	return Recover(signature, data)
}

func Signature(account *Account, data []byte) (ethereum.Signature, error) {
	msg := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))

	signature, err := account.wallet.SignDataWithPassphrase(*account.account, account.passphrase, "", msg)
	if err != nil {
		return ethereum.Signature{}, err
	}

	// Transform V from 0/1 to 27/28 according to the yellow paper:
	signature[64] += 27

	return ethereum.SignatureFromBytes(signature), nil
}

func Recover(signature ethereum.Signature, data []byte) (*ethereum.Address, error) {
	if signature[64] != 27 && signature[64] != 28 {
		return nil, ErrInvalidSignature
	}

	// Transform V from 27/28 to 0/1 according to yellow paper:
	signature[64] -= 27

	msg := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
	hash := crypto.Keccak256(msg)

	rpk, err := crypto.SigToPub(hash, signature[:])
	if err != nil {
		return nil, err
	}

	address := crypto.PubkeyToAddress(*rpk)
	return &address, nil
}
