package spire

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	datastoreMemory "github.com/toknowwhy/theunit-oracle/pkg/datastore/memory"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	ethereumMocks "github.com/toknowwhy/theunit-oracle/pkg/ethereum/mocks"
	"github.com/toknowwhy/theunit-oracle/pkg/log/null"
	"github.com/toknowwhy/theunit-oracle/pkg/spire"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/local"
)

func TestSpire_ConfigureAgent(t *testing.T) {
	prevSpireAgentFactory := spireAgentFactory
	defer func() {
		spireAgentFactory = prevSpireAgentFactory
	}()

	signer := &ethereumMocks.Signer{}
	transport := local.New(context.Background(), 0, nil)
	feeds := []ethereum.Address{ethereum.HexToAddress("0x07a35a1d4b751a818d93aa38e615c0df23064881")}
	logger := null.New()
	ds := &datastoreMemory.Datastore{}

	config := Spire{
		RPC:   RPC{Address: "1.2.3.4:1234"},
		Pairs: []string{"AAABBB"},
	}

	spireAgentFactory = func(ctx context.Context, cfg spire.AgentConfig) (*spire.Agent, error) {
		assert.NotNil(t, ctx)
		assert.Equal(t, ds, cfg.Datastore)
		assert.Equal(t, transport, cfg.Transport)
		assert.Equal(t, signer, cfg.Signer)
		assert.Equal(t, "tcp", cfg.Network)
		assert.Equal(t, "1.2.3.4:1234", cfg.Address)
		assert.Equal(t, logger, cfg.Logger)
		return &spire.Agent{}, nil
	}

	a, err := config.ConfigureAgent(AgentDependencies{
		Context:   context.Background(),
		Signer:    signer,
		Transport: transport,
		Datastore: ds,
		Feeds:     feeds,
		Logger:    logger,
	})
	require.NoError(t, err)
	require.NotNil(t, a)
}

func TestSpire_ConfigureClient(t *testing.T) {
	prevSpireClientFactory := spireClientFactory
	defer func() { spireClientFactory = prevSpireClientFactory }()

	signer := &ethereumMocks.Signer{}

	config := Spire{
		RPC:   RPC{Address: "1.2.3.4:1234"},
		Pairs: []string{"AAABBB"},
	}

	spireClientFactory = func(ctx context.Context, cfg spire.ClientConfig) (*spire.Client, error) {
		assert.NotNil(t, ctx)
		assert.Equal(t, signer, cfg.Signer)
		assert.Equal(t, "tcp", cfg.Network)
		assert.Equal(t, "1.2.3.4:1234", cfg.Address)
		return &spire.Client{}, nil
	}

	c, err := config.ConfigureClient(ClientDependencies{
		Context: context.Background(),
		Signer:  signer,
	})
	require.NoError(t, err)
	require.NotNil(t, c)
}
