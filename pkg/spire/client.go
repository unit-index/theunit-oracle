package spire

import (
	"context"
	"errors"
	"net/rpc"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/messages"
)

type Client struct {
	ctx    context.Context
	doneCh chan struct{}

	rpc     *rpc.Client
	network string
	address string
	signer  ethereum.Signer
}

type ClientConfig struct {
	Signer  ethereum.Signer
	Network string
	Address string
}

func NewClient(ctx context.Context, cfg ClientConfig) (*Client, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}
	return &Client{
		ctx:     ctx,
		doneCh:  make(chan struct{}),
		network: cfg.Network,
		address: cfg.Address,
		signer:  cfg.Signer,
	}, nil
}

func (c *Client) Start() error {
	client, err := rpc.DialHTTP(c.network, c.address)
	if err != nil {
		return err
	}
	c.rpc = client
	go c.contextCancelHandler()
	return nil
}

func (c *Client) Wait() {
	<-c.doneCh
}

func (c *Client) PublishPrice(price *messages.Price) error {
	err := c.rpc.Call("API.PublishPrice", PublishPriceArg{Price: price}, &Nothing{})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PullPrices(assetPair string, feeder string) ([]*messages.Price, error) {
	resp := &PullPricesResp{}
	err := c.rpc.Call("API.PullPrices", PullPricesArg{FilterAssetPair: assetPair, FilterFeeder: feeder}, resp)
	if err != nil {
		return nil, err
	}
	return resp.Prices, nil
}

func (c *Client) PullPrice(assetPair string, feeder string) (*messages.Price, error) {
	resp := &PullPriceResp{}
	err := c.rpc.Call("API.PullPrice", PullPriceArg{AssetPair: assetPair, Feeder: feeder}, resp)
	if err != nil {
		return nil, err
	}
	return resp.Price, nil
}

func (c *Client) contextCancelHandler() {
	defer func() { close(c.doneCh) }()
	<-c.ctx.Done()

	c.rpc.Close()
}
