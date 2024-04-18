package spire

import (
	"strings"

	"github.com/toknowwhy/theunit-oracle/pkg/datastore"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/messages"
)

type Nothing = struct{}

type API struct {
	transport transport.Transport
	datastore datastore.Datastore
	signer    ethereum.Signer
	log       log.Logger
}

type PublishPriceArg struct {
	Price *messages.Price
}

type PullPricesArg struct {
	FilterAssetPair string
	FilterFeeder    string
}

type PullPricesResp struct {
	Prices []*messages.Price
}

type PullPriceArg struct {
	AssetPair string
	Feeder    string
}

type PullPriceResp struct {
	Price *messages.Price
}

func (n *API) PublishPrice(arg *PublishPriceArg, _ *Nothing) error {
	n.log.
		WithFields(arg.Price.Price.Fields(n.signer)).
		Info("Publish price")

	return n.transport.Broadcast(messages.PriceMessageName, arg.Price)
}

func (n *API) PullPrices(arg *PullPricesArg, resp *PullPricesResp) error {
	n.log.
		WithField("assetPair", arg.FilterAssetPair).
		WithField("feeder", arg.FilterFeeder).
		Info("Pull prices")

	var prices []*messages.Price
	for fp, p := range n.datastore.Prices().All() {
		if arg.FilterAssetPair != "" && arg.FilterAssetPair != fp.AssetPair {
			continue
		}
		if arg.FilterFeeder != "" && !strings.EqualFold(arg.FilterFeeder, fp.Feeder.String()) {
			continue
		}
		prices = append(prices, p)
	}

	*resp = PullPricesResp{Prices: prices}

	return nil
}

func (n *API) PullPrice(arg *PullPriceArg, resp *PullPriceResp) error {
	n.log.
		WithField("assetPair", arg.AssetPair).
		WithField("feeder", arg.Feeder).
		Info("Pull price")

	*resp = PullPriceResp{
		Price: n.datastore.Prices().Feeder(arg.AssetPair, ethereum.HexToAddress(arg.Feeder)),
	}

	return nil
}
