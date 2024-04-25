package ghost

import (
	"testing"
)

func TestGhost_Configure(t *testing.T) {
	//prevGhostFactory := ghostFactory
	//defer func() { ghostFactory = prevGhostFactory }()
	//
	//interval := 10
	//pairs := []string{"AAABBB", "XXXYYY"}
	//gofer := &goferMocks.Gofer{}
	//signer := &ethereumMocks.Signer{}
	//transport := local.New(context.Background(), 0, nil)
	//logger := null.New()
	//
	//config := Ghost{
	//	Interval: interval,
	//	Pairs:    pairs,
	//}
	//
	//ghostFactory = func(ctx context.Context, cfg ghost.Config) (*ghost.Ghost, error) {
	//	assert.NotNil(t, ctx)
	//	assert.Equal(t, time.Duration(interval)*time.Second, cfg.Interval)
	//	assert.Equal(t, pairs, cfg.Pairs)
	//	assert.Equal(t, signer, cfg.Signer)
	//	assert.Equal(t, transport, cfg.Transport)
	//	assert.Equal(t, logger, cfg.Logger)
	//
	//	return &ghost.Ghost{}, nil
	//}

	//g, err := config.Configure(Dependencies{
	//	Context:   context.Background(),
	//	Gofer:     gofer,
	//	Signer:    signer,
	//	Transport: transport,
	//	Logger:    logger,
	//})
	//require.NoError(t, err)
	//assert.NotNil(t, g)
}
