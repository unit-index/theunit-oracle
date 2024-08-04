package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/toknowwhy/theunit-oracle/internal/config"
	ethereumConfig "github.com/toknowwhy/theunit-oracle/internal/config/ethereum"
	feedsConfig "github.com/toknowwhy/theunit-oracle/internal/config/feeds"
	ghostConfig "github.com/toknowwhy/theunit-oracle/internal/config/ghost"
	goferConfig "github.com/toknowwhy/theunit-oracle/internal/config/gofer"
	transportConfig "github.com/toknowwhy/theunit-oracle/internal/config/transport"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/ghost"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	logLogrus "github.com/toknowwhy/theunit-oracle/pkg/log/logrus"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"

	"github.com/toknowwhy/theunit-oracle/pkg/log"
)

type Config struct {
	Gofer     goferConfig.Gofer         `json:"gofer"`
	Ethereum  ethereumConfig.Ethereum   `json:"ethereum"`
	Transport transportConfig.Transport `json:"transport"`
	Ghost     ghostConfig.Ghost         `json:"ghost"`
	Feeds     feedsConfig.Feeds         `json:"feeds"`
}

type Dependencies struct {
	Context context.Context
	Logger  log.Logger
}

func (c *Config) Configure(d Dependencies, noGoferRPC bool) (transport.Transport, gofer.Gofer, *ghost.Ghost, error) {
	sig, err := c.Ethereum.ConfigureSigner()
	if err != nil {
		return nil, nil, nil, err
	}
	cli, err := c.Ethereum.ConfigureEthereumClient(nil) // signer may be empty here
	if err != nil {
		return nil, nil, nil, err
	}
	gof, err := c.Gofer.ConfigureGofer(d.Context, cli, d.Logger, noGoferRPC)
	if err != nil {
		return nil, nil, nil, err
	}

	if sig.Address() == ethereum.EmptyAddress {
		return nil, nil, nil, errors.New("ethereum account must be configured")
	}
	fed, err := c.Feeds.Addresses()
	if err != nil {
		return nil, nil, nil, err
	}
	tra, err := c.Transport.Configure(transportConfig.Dependencies{
		Context: d.Context,
		Signer:  sig,
		Feeds:   fed,
		Logger:  d.Logger,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	gho, err := c.Ghost.Configure(ghostConfig.Dependencies{
		Context:   d.Context,
		Gofer:     gof,
		Signer:    sig,
		Transport: tra,
		Logger:    d.Logger,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return tra, gof, gho, nil
}

type Services struct {
	ctxCancel context.CancelFunc
	Transport transport.Transport
	Gofer     gofer.Gofer
	Ghost     *ghost.Ghost
}

func PrepareServices(ctx context.Context, opts *options) (*Services, error) {
	var err error
	ctx, ctxCancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			ctxCancel()
		}
	}()

	// Load config file:
	err = config.ParseFile(&opts.Config, opts.ConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Logger:
	ll, err := logrus.ParseLevel(opts.LogVerbosity)
	if err != nil {
		return nil, err
	}
	lr := logrus.New()
	lr.SetLevel(ll)
	lr.SetFormatter(opts.LogFormat.Formatter())
	logger := logLogrus.New(lr)

	// Services:
	tra, gof, gho, err := opts.Config.Configure(Dependencies{
		Context: ctx,
		Logger:  logger,
	}, opts.GoferNoRPC)
	if err != nil {
		return nil, fmt.Errorf(" origin with namefailed to load Ghost configuration: %w", err)
	}

	return &Services{
		ctxCancel: ctxCancel,
		Transport: tra,
		Gofer:     gof,
		Ghost:     gho,
	}, nil
}

func (s *Services) Start() error {
	var err error
	if g, ok := s.Gofer.(gofer.StartableGofer); ok {
		if err = g.Start(); err != nil {
			return err
		}
	}
	if err = s.Transport.Start(); err != nil {
		return err
	}
	if err = s.Ghost.Start(); err != nil {
		return err
	}
	return nil
}

func (s *Services) CancelAndWait() {
	s.ctxCancel()
	s.Transport.Wait()
	s.Ghost.Wait()
	if g, ok := s.Gofer.(gofer.StartableGofer); ok {
		g.Wait()
	}
}
