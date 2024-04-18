package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

type Client struct {
	mock.Mock
}

func (c *Client) Call(ctx context.Context, call ethereum.Call) ([]byte, error) {
	args := c.Called(ctx, call)
	return args.Get(0).([]byte), args.Error(1)
}

func (c *Client) MultiCall(ctx context.Context, calls []ethereum.Call) ([][]byte, error) {
	args := c.Called(ctx, calls)
	return args.Get(0).([][]byte), args.Error(1)
}

func (c *Client) Storage(ctx context.Context, address ethereum.Address, key ethereum.Hash) ([]byte, error) {
	args := c.Called(ctx, address, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (c *Client) SendTransaction(ctx context.Context, transaction *ethereum.Transaction) (*ethereum.Hash, error) {
	args := c.Called(ctx, transaction)
	return args.Get(0).(*ethereum.Hash), args.Error(1)
}
