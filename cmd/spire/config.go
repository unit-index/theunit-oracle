package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/toknowwhy/theunit-oracle/internal/config"
	ethereumConfig "github.com/toknowwhy/theunit-oracle/internal/config/ethereum"
	feedsConfig "github.com/toknowwhy/theunit-oracle/internal/config/feeds"
	spireConfig "github.com/toknowwhy/theunit-oracle/internal/config/spire"
	transportConfig "github.com/toknowwhy/theunit-oracle/internal/config/transport"
	"github.com/toknowwhy/theunit-oracle/pkg/datastore"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	logLogrus "github.com/toknowwhy/theunit-oracle/pkg/log/logrus"
	"github.com/toknowwhy/theunit-oracle/pkg/spire"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
)

type Config struct {
	Transport transportConfig.Transport `json:"transport"`
	Ethereum  ethereumConfig.Ethereum   `json:"ethereum"`
	Spire     spireConfig.Spire         `json:"spire"`
	Feeds     feedsConfig.Feeds         `json:"feeds"`
}

type ClientDependencies struct {
	Context context.Context
}

type AgentDependencies struct {
	Context context.Context
	Logger  log.Logger
}

func (c *Config) ConfigureClient(d ClientDependencies) (*spire.Client, error) {
	sig, err := c.Ethereum.ConfigureSigner()
	if err != nil {
		return nil, err
	}
	cli, err := c.Spire.ConfigureClient(spireConfig.ClientDependencies{
		Context: d.Context,
		Signer:  sig,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (c *Config) ConfigureAgent(d AgentDependencies) (transport.Transport, datastore.Datastore, *spire.Agent, error) {
	sig, err := c.Ethereum.ConfigureSigner()
	if err != nil {
		return nil, nil, nil, err
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
	dat, err := c.Spire.ConfigureDatastore(spireConfig.DatastoreDependencies{
		Context:   d.Context,
		Signer:    sig,
		Transport: tra,
		Feeds:     fed,
		Logger:    d.Logger,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	age, err := c.Spire.ConfigureAgent(spireConfig.AgentDependencies{
		Context:   d.Context,
		Signer:    sig,
		Transport: tra,
		Datastore: dat,
		Feeds:     fed,
		Logger:    d.Logger,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return tra, dat, age, nil
}

type ClientServices struct {
	ctxCancel context.CancelFunc
	Client    *spire.Client
}

func PrepareClientServices(ctx context.Context, opts *options) (*ClientServices, error) {
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

	// Services:
	cli, err := opts.Config.ConfigureClient(ClientDependencies{
		Context: ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load Spire configuration: %w", err)
	}

	return &ClientServices{
		ctxCancel: ctxCancel,
		Client:    cli,
	}, nil
}

func (s *ClientServices) Start() error {
	var err error
	if err = s.Client.Start(); err != nil {
		return err
	}
	return nil
}

func (s *ClientServices) CancelAndWait() {
	s.ctxCancel()
	s.Client.Wait()
}

type AgentServices struct {
	ctxCancel context.CancelFunc
	Transport transport.Transport
	Datastore datastore.Datastore
	Agent     *spire.Agent
}

func PrepareAgentServices(ctx context.Context, opts *options) (*AgentServices, error) {
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
	tra, dat, age, err := opts.Config.ConfigureAgent(AgentDependencies{
		Context: ctx,
		Logger:  logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load Spire configuration: %w", err)
	}

	return &AgentServices{
		ctxCancel: ctxCancel,
		Transport: tra,
		Datastore: dat,
		Agent:     age,
	}, nil
}

func (s *AgentServices) Start() error {
	var err error
	if err = s.Transport.Start(); err != nil {
		return err
	}
	if err = s.Datastore.Start(); err != nil {
		return err
	}
	if err = s.Agent.Start(); err != nil {
		return err
	}
	return nil
}

func (s *AgentServices) CancelAndWait() {
	s.ctxCancel()
	s.Transport.Wait()
	s.Datastore.Wait()
	s.Agent.Wait()
}
