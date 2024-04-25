package main

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/toknowwhy/theunit-oracle/internal/config"
	transportConfig "github.com/toknowwhy/theunit-oracle/internal/config/transport"
	logLogrus "github.com/toknowwhy/theunit-oracle/pkg/log/logrus"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"

	"github.com/toknowwhy/theunit-oracle/pkg/log"
)

type Config struct {
	Transport transportConfig.Transport `json:"transport"`
}

type Dependencies struct {
	Context context.Context
	Logger  log.Logger
}

func (c *Config) Configure(d Dependencies) (transport.Transport, error) {
	tra, err := c.Transport.ConfigureP2PBoostrap(transportConfig.BootstrapDependencies{
		Context: d.Context,
		Logger:  d.Logger,
	})
	if err != nil {
		return nil, err
	}
	return tra, nil
}

type Service struct {
	ctxCancel context.CancelFunc
	Transport transport.Transport
}

func PrepareService(ctx context.Context, opts *options) (*Service, error) {
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
		return nil, err
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
	tra, err := opts.Config.Configure(Dependencies{
		Context: ctx,
		Logger:  logger,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		ctxCancel: ctxCancel,
		Transport: tra,
	}, nil
}

func (s *Service) Start() error {
	var err error
	if err = s.Transport.Start(); err != nil {
		return err
	}
	return nil
}

func (s *Service) CancelAndWait() {
	s.ctxCancel()
	s.Transport.Wait()
}
