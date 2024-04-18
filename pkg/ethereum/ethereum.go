package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Aliases for the go-ethereum types and functions used in multiple packages.
// These aliases was created to not rely directly on the go-ethereum packages.

// AddressLength is the expected length of the address
const AddressLength = common.AddressLength

type (
	Address = common.Address
	Hash    = common.Hash
)

// HexToAddress returns Address from hex representation.
var HexToAddress = common.HexToAddress

// IsHexAddress verifies if given string is a valid Ethereum address.
var IsHexAddress = common.IsHexAddress

// EmptyAddress contains empty Ethereum address: 0x0000000000000000000000000000000000000000
var EmptyAddress Address

// HexToBytes returns bytes from hex string.
var HexToBytes = common.FromHex

// SHA3Hash calculates SHA3 hash.
func SHA3Hash(b []byte) []byte {
	return crypto.Keccak256Hash(b).Bytes()
}
