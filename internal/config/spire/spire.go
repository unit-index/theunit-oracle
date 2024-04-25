package spire

import (
	"context"

	"github.com/toknowwhy/theunit-oracle/pkg/datastore"
	datastoreMemory "github.com/toknowwhy/theunit-oracle/pkg/datastore/memory"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/spire"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
)

// nolint
var spireAgentFactory = func(ctx context.Context, cfg spire.AgentConfig) (*spire.Agent, error) {
	return spire.NewAgent(ctx, cfg)
}

// nolint
var spireClientFactory = func(ctx context.Context, cfg spire.ClientConfig) (*spire.Client, error) {
	return spire.NewClient(ctx, cfg)
}

var datastoreFactory = func(ctx context.Context, cfg datastoreMemory.Config) (datastore.Datastore, error) {
	return datastoreMemory.NewDatastore(ctx, cfg)
}

type Spire struct {
	RPC   RPC      `json:"rpc"`
	Pairs []string `json:"pairs"`
}

type RPC struct {
	Address string `json:"address"`
}

type AgentDependencies struct {
	Context   context.Context
	Signer    ethereum.Signer
	Transport transport.Transport
	Datastore datastore.Datastore
	Feeds     []ethereum.Address
	Logger    log.Logger
}

type ClientDependencies struct {
	Context context.Context
	Signer  ethereum.Signer
}

type DatastoreDependencies struct {
	Context   context.Context
	Signer    ethereum.Signer
	Transport transport.Transport
	Feeds     []ethereum.Address
	Logger    log.Logger
}

func (c *Spire) ConfigureAgent(d AgentDependencies) (*spire.Agent, error) {
	agent, err := spireAgentFactory(d.Context, spire.AgentConfig{
		Datastore: d.Datastore,
		Transport: d.Transport,
		Signer:    d.Signer,
		Network:   "tcp",
		Address:   c.RPC.Address,
		Logger:    d.Logger,
	})
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (c *Spire) ConfigureClient(d ClientDependencies) (*spire.Client, error) {
	return spireClientFactory(d.Context, spire.ClientConfig{
		Signer:  d.Signer,
		Network: "tcp",
		Address: c.RPC.Address,
	})
}

func (c *Spire) ConfigureDatastore(d DatastoreDependencies) (datastore.Datastore, error) {
	cfg := datastoreMemory.Config{
		Signer:    d.Signer,
		Transport: d.Transport,
		Pairs:     make(map[string]*datastoreMemory.Pair),
		Logger:    d.Logger,
	}
	for _, name := range c.Pairs {
		cfg.Pairs[name] = &datastoreMemory.Pair{Feeds: d.Feeds}
	}
	return datastoreFactory(d.Context, cfg)
}
