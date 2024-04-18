package spire

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/rpc"

	"github.com/toknowwhy/theunit-oracle/pkg/datastore"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
)

const AgentLoggerTag = "SPIRE_AGENT"

type Agent struct {
	ctx    context.Context
	doneCh chan struct{}

	api      *API
	rpc      *rpc.Server
	listener net.Listener
	network  string
	address  string
	log      log.Logger
}

type AgentConfig struct {
	Datastore datastore.Datastore
	Transport transport.Transport
	Signer    ethereum.Signer
	Network   string
	Address   string
	Logger    log.Logger
}

func NewAgent(ctx context.Context, cfg AgentConfig) (*Agent, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}
	s := &Agent{
		ctx:    ctx,
		doneCh: make(chan struct{}),
		api: &API{
			datastore: cfg.Datastore,
			transport: cfg.Transport,
			signer:    cfg.Signer,
			log:       cfg.Logger.WithField("tag", AgentLoggerTag),
		},
		rpc:     rpc.NewServer(),
		network: cfg.Network,
		address: cfg.Address,
		log:     cfg.Logger.WithField("tag", AgentLoggerTag),
	}
	err := s.rpc.Register(s.api)
	if err != nil {
		return nil, err
	}
	s.rpc.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	return s, nil
}

func (s *Agent) Start() error {
	s.log.Infof("Starting")
	var err error

	// Start RPC server:
	s.listener, err = net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	go func() {
		err := http.Serve(s.listener, nil)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.WithError(err).Error("RPC server crashed")
		}
	}()

	go s.contextCancelHandler()
	return nil
}

// Wait waits until agent's context is cancelled.
func (s *Agent) Wait() {
	<-s.doneCh
}

func (s *Agent) contextCancelHandler() {
	defer func() { close(s.doneCh) }()
	defer s.log.Info("Stopped")
	<-s.ctx.Done()

	err := s.listener.Close()
	if err != nil {
		s.log.WithError(err).Error("Unable to close RPC listener")
	}
}
