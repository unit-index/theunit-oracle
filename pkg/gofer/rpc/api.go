package rpc

import (
	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/graph/feeder"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
)

type Nothing = struct{}

type API struct {
	gofer gofer.Gofer
	log   log.Logger
}

type FeedArg struct {
	Pairs []gofer.Pair
}

type FeedResp struct {
	Warnings feeder.Warnings
}

type NodesArg struct {
	Format marshal.FormatType
	Pairs  []gofer.Pair
}

type NodesResp struct {
	Pairs map[gofer.Pair]*gofer.Model
}

type PricesArg struct {
	Pairs []gofer.Pair
}

type PricesResp struct {
	Prices map[gofer.Pair]*gofer.Price
}

type PairsResp struct {
	Pairs []gofer.Pair
}

func (n *API) Models(arg *NodesArg, resp *NodesResp) error {
	n.log.WithField("pairs", arg.Pairs).Info("Models")
	pairs, err := n.gofer.Models(arg.Pairs...)
	if err != nil {
		return err
	}
	resp.Pairs = pairs
	return nil
}

func (n *API) Prices(arg *PricesArg, resp *PricesResp) error {
	n.log.WithField("pairs", arg.Pairs).Info("Prices")
	prices, err := n.gofer.Prices(arg.Pairs...)
	if err != nil {
		return err
	}
	resp.Prices = prices
	return nil
}

func (n *API) Pairs(_ *Nothing, resp *PairsResp) error {
	n.log.Info("Prices")
	pairs, err := n.gofer.Pairs()
	if err != nil {
		return err
	}
	resp.Pairs = pairs
	return nil
}
