package local

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/pkg/transport"
)

type testMsg struct {
	Val string
}

func (t *testMsg) Marshall() ([]byte, error) {
	return []byte(t.Val), nil
}

func (t *testMsg) Unmarshall(bytes []byte) error {
	t.Val = string(bytes)
	return nil
}

func TestLocal_Broadcast(t *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	l := New(ctx, 1, map[string]transport.Message{"foo": (*testMsg)(nil)})

	// Valid message:
	vm := &testMsg{Val: "bar"}
	assert.NoError(t, l.Broadcast("foo", vm))
}

func TestLocal_Messages(t *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	l := New(ctx, 1, map[string]transport.Message{"foo": (*testMsg)(nil)})

	// Valid message:
	assert.NoError(t, l.Broadcast("foo", &testMsg{Val: "bar"}))
	assert.Equal(t, &testMsg{Val: "bar"}, (<-l.Messages("foo")).Message)
}
