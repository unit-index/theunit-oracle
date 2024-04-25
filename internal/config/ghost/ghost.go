package ghost

import (
	"context"
	"time"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/ghost"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
)

// nolint
var ghostFactory = func(ctx context.Context, cfg ghost.Config) (*ghost.Ghost, error) {
	return ghost.NewGhost(ctx, cfg)
}

type Ghost struct {
	Interval int      `json:"interval"`
	Pairs    []string `json:"pairs"`
}

type Dependencies struct {
	Context   context.Context
	Gofer     gofer.Gofer
	Signer    ethereum.Signer
	Transport transport.Transport
	Logger    log.Logger
}

func (c *Ghost) Configure(d Dependencies) (*ghost.Ghost, error) {
	cfg := ghost.Config{
		Gofer:     d.Gofer,
		Signer:    d.Signer,
		Transport: d.Transport,
		Logger:    d.Logger,
		Interval:  time.Second * time.Duration(c.Interval),
		Pairs:     c.Pairs,
	}
	return ghostFactory(d.Context, cfg)
}
