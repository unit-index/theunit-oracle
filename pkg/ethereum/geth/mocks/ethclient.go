package mocks

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

type EthClient struct {
	mock.Mock
}

func (e *EthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	args := e.Called(ctx, tx)
	return args.Error(0)
}

func (e *EthClient) StorageAt(ctx context.Context, acc common.Address, key common.Hash, block *big.Int) ([]byte, error) {
	args := e.Called(ctx, acc, key, block)
	return args.Get(0).([]byte), args.Error(1)
}

func (e *EthClient) CallContract(ctx context.Context, call ethereum.CallMsg, block *big.Int) ([]byte, error) {
	args := e.Called(ctx, call, block)
	return args.Get(0).([]byte), args.Error(1)
}

func (e *EthClient) NonceAt(ctx context.Context, account common.Address, block *big.Int) (uint64, error) {
	args := e.Called(ctx, account, block)
	return uint64(args.Int(0)), args.Error(1)
}

func (e *EthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	args := e.Called(ctx, account)
	return uint64(args.Int(0)), args.Error(1)
}

func (e *EthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := e.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (e *EthClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	args := e.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (e *EthClient) NetworkID(ctx context.Context) (*big.Int, error) {
	args := e.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}
