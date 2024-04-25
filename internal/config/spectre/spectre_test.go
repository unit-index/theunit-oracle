package spectre

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	datastoreMemory "github.com/toknowwhy/theunit-oracle/pkg/datastore/memory"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	ethereumMocks "github.com/toknowwhy/theunit-oracle/pkg/ethereum/mocks"
	"github.com/toknowwhy/theunit-oracle/pkg/log/null"
	"github.com/toknowwhy/theunit-oracle/pkg/spectre"
)

func TestSpectre_Configure(t *testing.T) {
	prevSpectreFactory := spectreFactory
	prevDatastoreFactory := datastoreFactory
	defer func() {
		spectreFactory = prevSpectreFactory
		datastoreFactory = prevDatastoreFactory
	}()

	interval := int64(10)
	signer := &ethereumMocks.Signer{}
	ethClient := &ethereumMocks.Client{}
	feeds := []ethereum.Address{ethereum.HexToAddress("0x07a35a1d4b751a818d93aa38e615c0df23064881")}
	ds := &datastoreMemory.Datastore{}
	logger := null.New()

	config := Spectre{
		Interval: interval,
		Medianizers: map[string]Medianizer{
			"AAABBB": {
				Contract:         "0xe0F30cb149fAADC7247E953746Be9BbBB6B5751f",
				OracleSpread:     0.1,
				OracleExpiration: 15500,
				MsgExpiration:    1800,
			},
		},
	}

	spectreFactory = func(ctx context.Context, cfg spectre.Config) (*spectre.Spectre, error) {
		assert.NotNil(t, ctx)
		assert.Equal(t, signer, cfg.Signer)
		assert.Equal(t, ds, cfg.Datastore)
		assert.Equal(t, secToDuration(interval), cfg.Interval)
		assert.Equal(t, logger, cfg.Logger)
		assert.Equal(t, "AAABBB", cfg.Pairs[0].AssetPair)
		assert.Equal(t, secToDuration(config.Medianizers["AAABBB"].OracleExpiration), cfg.Pairs[0].OracleExpiration)
		assert.Equal(t, secToDuration(config.Medianizers["AAABBB"].MsgExpiration), cfg.Pairs[0].PriceExpiration)
		assert.Equal(t, config.Medianizers["AAABBB"].OracleSpread, cfg.Pairs[0].OracleSpread)
		assert.Equal(t, ethereum.HexToAddress(config.Medianizers["AAABBB"].Contract), cfg.Pairs[0].Median.Address())
		return &spectre.Spectre{}, nil
	}

	s, err := config.ConfigureSpectre(Dependencies{
		Context:        context.Background(),
		Signer:         signer,
		Datastore:      ds,
		EthereumClient: ethClient,
		Feeds:          feeds,
		Logger:         logger,
	})
	require.NoError(t, err)
	require.NotNil(t, s)
}

func secToDuration(s int64) time.Duration {
	return time.Duration(s) * time.Second
}
