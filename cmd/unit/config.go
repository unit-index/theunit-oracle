package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/toknowwhy/theunit-oracle/internal/config"
	ethereumConfig "github.com/toknowwhy/theunit-oracle/internal/config/ethereum"
	feedsConfig "github.com/toknowwhy/theunit-oracle/internal/config/feeds"
	goferConfig "github.com/toknowwhy/theunit-oracle/internal/config/gofer"
	transportConfig "github.com/toknowwhy/theunit-oracle/internal/config/transport"
	unitConfig "github.com/toknowwhy/theunit-oracle/internal/config/unit"
	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	logLogrus "github.com/toknowwhy/theunit-oracle/pkg/log/logrus"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
	pkgUnit "github.com/toknowwhy/theunit-oracle/pkg/unit"
)

type Config struct {
	Ethereum  ethereumConfig.Ethereum   `json:"ethereum"`
	Unit      unitConfig.Unit           `json:"unit"`
	Gofer     goferConfig.Gofer         `json:"gofer"`
	Transport transportConfig.Transport `json:"transport"`
	Feeds     feedsConfig.Feeds         `json:"feeds"`
}

type Dependencies struct {
	Context context.Context
	Logger  log.Logger
}

func (c *Config) Configure(ctx context.Context, logger log.Logger, noRPC bool) (pkgUnit.Unit, error) {
	cli, err := c.Ethereum.ConfigureEthereumClient(nil)
	if err != nil {
		return nil, err
	}
	gfo, err := c.Gofer.ConfigureGofer(ctx, cli, logger, noRPC)
	if err != nil {
		return nil, err
	}
	signer, err := c.Ethereum.ConfigureSigner()
	if err != nil {
		return nil, err
	}
	fed, err := c.Feeds.Addresses()
	if err != nil {
		return nil, err
	}

	transport, err := c.Transport.ConfigureUnit(transportConfig.Dependencies{
		Context: ctx,
		Signer:  signer,
		Feeds:   fed,
		Logger:  logger,
	})
	return c.Unit.ConfigureUnit(ctx, cli, gfo, logger, noRPC, signer, transport, fed)
}

func (c *Config) ConfigureFeed(d Dependencies, noGoferRPC bool) (transport.Transport, gofer.Gofer, pkgUnit.Unit, error) {
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
	//tra, err := c.Transport.Configure(transportConfig.Dependencies{
	//	Context: d.Context,
	//	Signer:  sig,
	//	Feeds:   fed,
	//	Logger:  d.Logger,
	//})
	//if err != nil {
	//	return nil, nil, nil, err
	//}
	gfo, err := c.Gofer.ConfigureGofer(d.Context, cli, d.Logger, noGoferRPC)
	if err != nil {
		return nil, nil, nil, err
	}

	signer, err := c.Ethereum.ConfigureSigner()
	if err != nil {
		return nil, nil, nil, err
	}

	transport, err := c.Transport.ConfigureUnit(transportConfig.Dependencies{
		Context: d.Context,
		Signer:  signer,
		Feeds:   fed,
		Logger:  d.Logger,
	})

	unit, err := c.Unit.ConfigureUnit(d.Context, cli, gfo, d.Logger, noGoferRPC, signer, transport, fed)
	if err != nil {
		return nil, nil, nil, err
	}
	return transport, gof, unit, nil
}

//func (c *Config) ConfigureRPCAgent(ctx context.Context, logger log.Logger) (*rpc.Agent, error) {
//	cli, err := c.Ethereum.ConfigureEthereumClient(nil)
//	if err != nil {
//		return nil, err
//	}
//	return c.Gofer.ConfigureRPCAgent(ctx, cli, logger)
//}

type UnitClientServices struct {
	ctxCancel  context.CancelFunc
	Unit       pkgUnit.Unit
	Marshaller marshal.Marshaller
}

type UnitServerServices struct {
	ctxCancel  context.CancelFunc
	Unit       pkgUnit.Unit
	Marshaller marshal.Marshaller
	Transport  transport.Transport
	Gofer      gofer.Gofer
}

func PrepareUnitClientServices(ctx context.Context, opts *options) (*UnitClientServices, error) {
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
	unit, err := opts.Config.Configure(ctx, logger, opts.NoRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to load Gofer configuration1: %w", err)
	}
	mar, err := marshal.NewMarshal(opts.Format.format)
	if err != nil {
		return nil, err
	}

	return &UnitClientServices{
		ctxCancel:  ctxCancel,
		Unit:       unit,
		Marshaller: mar,
	}, nil
}

func PrepareUnitServerServices(ctx context.Context, opts *options) (*UnitServerServices, error) {
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
	tra, gof, unit, err := opts.Config.ConfigureFeed(Dependencies{Context: ctx, Logger: logger}, opts.NoRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to load Gofer configuration1: %w", err)
	}
	mar, err := marshal.NewMarshal(opts.Format.format)
	if err != nil {
		return nil, err
	}

	return &UnitServerServices{
		ctxCancel:  ctxCancel,
		Unit:       unit,
		Gofer:      gof,
		Transport:  tra,
		Marshaller: mar,
	}, nil
}

func (s *UnitClientServices) Start() error {
	if g, ok := s.Unit.(pkgUnit.StartableUnit); ok {
		return g.Start()
	}
	return nil
}

func (s *UnitServerServices) Start() error {
	if g, ok := s.Gofer.(gofer.StartableGofer); ok {
		if err := g.Start(); err != nil {
			return err
		}
	}
	if err := s.Transport.Start(); err != nil {
		fmt.Println(err)
		return err
	}

	if u, ok := s.Unit.(pkgUnit.StartableUnit); ok {
		return u.Start()
	}
	return nil
}

func (s *UnitClientServices) CancelAndWait() {
	s.ctxCancel()
	if g, ok := s.Unit.(pkgUnit.StartableUnit); ok {
		g.Wait()
	}
}

func (s *UnitServerServices) CancelAndWait() {
	s.ctxCancel()
	if g, ok := s.Unit.(pkgUnit.StartableUnit); ok {
		g.Wait()
	}
}

//type GoferAgentService struct {
//	ctxCancel context.CancelFunc
//	Agent     *rpc.Agent
//}

//func PrepareGoferAgentService(ctx context.Context, opts *options) (*GoferAgentService, error) {
//	var err error
//	ctx, ctxCancel := context.WithCancel(ctx)
//	defer func() {
//		if err != nil {
//			ctxCancel()
//		}
//	}()
//
//	// Load config file:
//	err = config.ParseFile(&opts.Config, opts.ConfigFilePath)
//	if err != nil {
//		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
//	}
//
//	// Logger:
//	ll, err := logrus.ParseLevel(opts.LogVerbosity)
//	if err != nil {
//		return nil, err
//	}
//	lr := logrus.New()
//	lr.SetLevel(ll)
//	lr.SetFormatter(opts.LogFormat.Formatter())
//	logger := logLogrus.New(lr)
//
//	// Services:
//	age, err := opts.Config.ConfigureRPCAgent(ctx, logger)
//	if err != nil {
//		return nil, fmt.Errorf("failed to load Gofer configuration: %w", err)
//	}
//
//	return &GoferAgentService{
//		ctxCancel: ctxCancel,
//		Agent:     age,
//	}, nil
//}

//func (s *GoferAgentService) Start() error {
//	return s.Agent.Start()
//}
//
//func (s *GoferAgentService) CancelAndWait() {
//	s.ctxCancel()
//	s.Agent.Wait()
//}
