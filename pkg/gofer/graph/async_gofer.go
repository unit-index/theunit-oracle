package graph

import (
	"context"
	"errors"

	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/graph/feeder"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/graph/nodes"
)

// AsyncGofer implements the gofer.Gofer interface. It works just like Graph
// but allows to update prices asynchronously.
type AsyncGofer struct {
	*Gofer
	ctx    context.Context
	feeder *feeder.Feeder
	doneCh chan struct{}
}

// NewAsyncGofer returns a new AsyncGofer instance.
func NewAsyncGofer(ctx context.Context, g map[gofer.Pair]nodes.Aggregator, f *feeder.Feeder) (*AsyncGofer, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}
	return &AsyncGofer{
		Gofer:  NewGofer(g, nil),
		ctx:    ctx,
		feeder: f,
		doneCh: make(chan struct{}),
	}, nil
}

// Start starts asynchronous price updater.
func (a *AsyncGofer) Start() error {
	go a.contextCancelHandler()
	ns, _ := a.findNodes()
	return a.feeder.Start(ns...)
}

// Wait waits until feeder's context is cancelled.
func (a *AsyncGofer) Wait() {
	<-a.doneCh
}

func (a *AsyncGofer) contextCancelHandler() {
	defer func() { close(a.doneCh) }()
	<-a.ctx.Done()

	a.feeder.Wait()
}
