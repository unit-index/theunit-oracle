package rpc

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/mocks"
	"github.com/toknowwhy/theunit-oracle/pkg/log/null"
)

var (
	agent     *Agent
	mockGofer *mocks.Gofer
	rpcGofer  *Gofer
)

func TestMain(m *testing.M) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	var err error

	mockGofer = &mocks.Gofer{}
	agent, err = NewAgent(ctx, AgentConfig{
		Gofer:   mockGofer,
		Network: "tcp",
		Address: "127.0.0.1:0",
		Logger:  null.New(),
	})
	if err != nil {
		panic(err)
	}
	if err = agent.Start(); err != nil {
		panic(err)
	}
	rpcGofer, err = NewGofer(ctx, "tcp", agent.listener.Addr().String())
	if err != nil {
		panic(err)
	}
	err = rpcGofer.Start()
	if err != nil {
		panic(err)
	}

	retCode := m.Run()
	ctxCancel()
	os.Exit(retCode)
}

func TestClient_Models(t *testing.T) {
	pair := gofer.Pair{Base: "A", Quote: "B"}
	model := map[gofer.Pair]*gofer.Model{pair: {Type: "test"}}

	mockGofer.On("Models", pair).Return(model, nil)
	resp, err := rpcGofer.Models(pair)

	assert.Equal(t, model, resp)
	assert.NoError(t, err)
}

func TestClient_Price(t *testing.T) {
	pair := gofer.Pair{Base: "A", Quote: "B"}
	prices := map[gofer.Pair]*gofer.Price{pair: {Type: "test"}}

	mockGofer.On("Prices", pair).Return(prices, nil)
	resp, err := rpcGofer.Price(pair)

	assert.Equal(t, prices[pair], resp)
	assert.NoError(t, err)
}

func TestClient_Prices(t *testing.T) {
	pair := gofer.Pair{Base: "A", Quote: "B"}
	prices := map[gofer.Pair]*gofer.Price{pair: {Type: "test"}}

	mockGofer.On("Prices", pair).Return(prices, nil)
	resp, err := rpcGofer.Prices(pair)

	assert.Equal(t, prices, resp)
	assert.NoError(t, err)
}

func TestClient_Pairs(t *testing.T) {
	pairs := []gofer.Pair{{Base: "A", Quote: "B"}}

	mockGofer.On("Pairs").Return(pairs, nil)
	resp, err := rpcGofer.Pairs()

	assert.Equal(t, pairs, resp)
	assert.NoError(t, err)
}
