package ethereum

import (
	"context"
	"math/big"
)

type Transaction struct {
	// Address is the contract's address.
	Address Address
	// Nonce is the transaction nonce. If zero, the nonce will be filled
	// automatically.
	Nonce uint64
	// PriorityFee is the maximum tip value. If nil, the suggested gas tip value
	// will be used.
	PriorityFee *big.Int
	// MaxFee is the maximum fee value. If nil, double value of a suggested
	// gas fee will be used.
	MaxFee *big.Int
	// GasLimit is the maximum gas available to be used for this transaction.
	GasLimit *big.Int
	// Data is the raw transaction data.
	Data []byte
	// ChainID is the transaction chain ID. If nil, the chan ID will be filled
	// automatically.
	ChainID *big.Int
	// SignedTx contains signed transaction. The data type stored here may
	// be different for various implementations.
	SignedTx interface{}
}

type Call struct {
	// Address is the contract's address.
	Address Address
	// Data is the raw call data.
	Data []byte
}

type Client interface {
	// Call executes a message call transaction, which is directly
	// executed in the VM of the node, but never mined into the blockchain.
	Call(ctx context.Context, call Call) ([]byte, error)
	// MultiCall works like the Call function but allows to execute multiple
	// calls at once.
	MultiCall(ctx context.Context, calls []Call) ([][]byte, error)
	// Storage returns the value of key in the contract storage of the
	// given account.
	Storage(ctx context.Context, address Address, key Hash) ([]byte, error)
	// SendTransaction injects a signed transaction into the pending pool
	// for execution.
	SendTransaction(ctx context.Context, transaction *Transaction) (*Hash, error)
}
