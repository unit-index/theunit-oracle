package oracle

import (
	"context"
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"math/big"
	"time"
)

type UnitAlgorithm interface {
	GetTokens(ctx context.Context, time *big.Int) ([]common.Address, error)
}

type UnitParams struct {
	Name          string
	LastMarketCap *big.Int
	LastPrice     *big.Int
	Age           time.Time
	// Signature:
	V byte
	R [32]byte
	S [32]byte

	// StarkWare signature:
	StarkR  []byte
	StarkS  []byte
	StarkPK []byte
}

func (u *UnitParams) Sign(signer ethereum.Signer) error {
	if u.LastMarketCap == nil {
		return ErrPriceNotSet
	}
	if u.LastPrice == nil {
		return ErrPriceNotSet
	}

	signature, err := signer.Signature(u.hash())
	if err != nil {
		return err
	}

	u.V, u.R, u.S = signature.VRS()

	return nil
}

// hash is an equivalent of keccak256(abi.encodePacked(val_, age_, wat))) in Solidity.
func (u *UnitParams) hash() []byte {

	Name := make([]byte, 32)
	copy(Name, u.Name)

	LastMarketCap := make([]byte, 32)
	u.LastMarketCap.FillBytes(LastMarketCap)

	// Time:
	age := make([]byte, 32)
	binary.BigEndian.PutUint64(age[24:], uint64(u.Age.Unix()))

	// Asset name:
	LastPrice := make([]byte, 32)
	u.LastPrice.FillBytes(LastPrice)

	hash := make([]byte, 128)
	copy(hash[0:32], Name)
	copy(hash[32:64], age)
	copy(hash[64:96], LastMarketCap)
	copy(hash[96:128], LastPrice)

	return ethereum.SHA3Hash(hash)
}

func (u *UnitParams) From(signer ethereum.Signer) (*ethereum.Address, error) {
	from, err := signer.Recover(u.Signature(), u.hash())
	if err != nil {
		return nil, err
	}

	return from, nil
}

func (u *UnitParams) Signature() ethereum.Signature {
	return ethereum.SignatureFromVRS(u.V, u.R, u.S)
}
