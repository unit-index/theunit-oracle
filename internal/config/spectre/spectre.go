package spectre

import (
	"context"
	"time"

	"github.com/toknowwhy/theunit-oracle/pkg/datastore"
	datastoreMemory "github.com/toknowwhy/theunit-oracle/pkg/datastore/memory"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	oracleGeth "github.com/toknowwhy/theunit-oracle/pkg/oracle/geth"
	"github.com/toknowwhy/theunit-oracle/pkg/spectre"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
)

// nolint
var spectreFactory = func(ctx context.Context, cfg spectre.Config) (*spectre.Spectre, error) {
	return spectre.NewSpectre(ctx, cfg)
}

var datastoreFactory = func(ctx context.Context, cfg datastoreMemory.Config) (datastore.Datastore, error) {
	return datastoreMemory.NewDatastore(ctx, cfg)
}

type Spectre struct {
	Interval    int64                 `json:"interval"`
	Medianizers map[string]Medianizer `json:"medianizers"`
}

type Medianizer struct {
	Contract         string  `json:"oracle"`
	OracleSpread     float64 `json:"oracleSpread"`
	OracleExpiration int64   `json:"oracleExpiration"`
	MsgExpiration    int64   `json:"msgExpiration"`
}

type Dependencies struct {
	Context        context.Context
	Signer         ethereum.Signer
	Datastore      datastore.Datastore
	EthereumClient ethereum.Client
	Feeds          []ethereum.Address
	Logger         log.Logger
}

type DatastoreDependencies struct {
	Context   context.Context
	Signer    ethereum.Signer
	Transport transport.Transport
	Feeds     []ethereum.Address
	Logger    log.Logger
}

func (c *Spectre) ConfigureSpectre(d Dependencies) (*spectre.Spectre, error) {
	cfg := spectre.Config{
		Signer:    d.Signer,
		Interval:  time.Second * time.Duration(c.Interval),
		Datastore: d.Datastore,
		Logger:    d.Logger,
	}
	for name, pair := range c.Medianizers {
		cfg.Pairs = append(cfg.Pairs, &spectre.Pair{
			AssetPair:        name,
			OracleSpread:     pair.OracleSpread,
			OracleExpiration: time.Second * time.Duration(pair.OracleExpiration),
			PriceExpiration:  time.Second * time.Duration(pair.MsgExpiration),
			Median:           oracleGeth.NewMedian(d.EthereumClient, ethereum.HexToAddress(pair.Contract)),
		})
	}
	return spectreFactory(d.Context, cfg)
}

func (c *Spectre) ConfigureDatastore(d DatastoreDependencies) (datastore.Datastore, error) {
	cfg := datastoreMemory.Config{
		Signer:    d.Signer,
		Transport: d.Transport,
		Pairs:     make(map[string]*datastoreMemory.Pair),
		Logger:    d.Logger,
	}
	for name := range c.Medianizers {
		cfg.Pairs[name] = &datastoreMemory.Pair{Feeds: d.Feeds}
	}
	return datastoreFactory(d.Context, cfg)
}
