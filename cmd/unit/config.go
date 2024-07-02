package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	pkgUnit "github.com/toknowwhy/theunit-oracle/pkg/unit"

	"github.com/toknowwhy/theunit-oracle/internal/config"
	ethereumConfig "github.com/toknowwhy/theunit-oracle/internal/config/ethereum"
	unitConfig "github.com/toknowwhy/theunit-oracle/internal/config/unit"
	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	logLogrus "github.com/toknowwhy/theunit-oracle/pkg/log/logrus"
)

type Config struct {
	Ethereum ethereumConfig.Ethereum `json:"ethereum"`
	Unit     unitConfig.Unit         `json:"unit"`
}

func (c *Config) Configure(ctx context.Context, logger log.Logger, noRPC bool) (pkgUnit.Unit, error) {
	cli, err := c.Ethereum.ConfigureEthereumClient(nil)
	if err != nil {
		return nil, err
	}
	return c.Unit.ConfigureUnit(ctx, cli, logger, noRPC)
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
	gof, err := opts.Config.Configure(ctx, logger, opts.NoRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to load Gofer configuration1: %w", err)
	}
	mar, err := marshal.NewMarshal(opts.Format.format)
	if err != nil {
		return nil, err
	}

	return &UnitClientServices{
		ctxCancel:  ctxCancel,
		Unit:       gof,
		Marshaller: mar,
	}, nil
}

func (s *UnitClientServices) Start() error {
	if g, ok := s.Unit.(pkgUnit.StartableUnit); ok {
		return g.Start()
	}
	return nil
}

func (s *UnitClientServices) CancelAndWait() {
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
