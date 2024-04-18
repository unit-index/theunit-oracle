package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/toknowwhy/theunit-oracle/internal/config"
	ethereumConfig "github.com/toknowwhy/theunit-oracle/internal/config/ethereum"
	goferConfig "github.com/toknowwhy/theunit-oracle/internal/config/gofer"
	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal"
	pkgGofer "github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/rpc"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	logLogrus "github.com/toknowwhy/theunit-oracle/pkg/log/logrus"
)

type Config struct {
	Ethereum ethereumConfig.Ethereum `json:"ethereum"`
	Gofer    goferConfig.Gofer       `json:"gofer"`
}

func (c *Config) Configure(ctx context.Context, logger log.Logger, noRPC bool) (pkgGofer.Gofer, error) {
	cli, err := c.Ethereum.ConfigureEthereumClient(nil)
	if err != nil {
		return nil, err
	}
	return c.Gofer.ConfigureGofer(ctx, cli, logger, noRPC)
}

func (c *Config) ConfigureRPCAgent(ctx context.Context, logger log.Logger) (*rpc.Agent, error) {
	cli, err := c.Ethereum.ConfigureEthereumClient(nil)
	if err != nil {
		return nil, err
	}
	return c.Gofer.ConfigureRPCAgent(ctx, cli, logger)
}

type GoferClientServices struct {
	ctxCancel  context.CancelFunc
	Gofer      pkgGofer.Gofer
	Marshaller marshal.Marshaller
}

func PrepareGoferClientServices(ctx context.Context, opts *options) (*GoferClientServices, error) {
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

	return &GoferClientServices{
		ctxCancel:  ctxCancel,
		Gofer:      gof,
		Marshaller: mar,
	}, nil
}

func (s *GoferClientServices) Start() error {
	if g, ok := s.Gofer.(pkgGofer.StartableGofer); ok {
		return g.Start()
	}
	return nil
}

func (s *GoferClientServices) CancelAndWait() {
	s.ctxCancel()
	if g, ok := s.Gofer.(pkgGofer.StartableGofer); ok {
		g.Wait()
	}
}

type GoferAgentService struct {
	ctxCancel context.CancelFunc
	Agent     *rpc.Agent
}

func PrepareGoferAgentService(ctx context.Context, opts *options) (*GoferAgentService, error) {
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
	age, err := opts.Config.ConfigureRPCAgent(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to load Gofer configuration: %w", err)
	}

	return &GoferAgentService{
		ctxCancel: ctxCancel,
		Agent:     age,
	}, nil
}

func (s *GoferAgentService) Start() error {
	return s.Agent.Start()
}

func (s *GoferAgentService) CancelAndWait() {
	s.ctxCancel()
	s.Agent.Wait()
}
