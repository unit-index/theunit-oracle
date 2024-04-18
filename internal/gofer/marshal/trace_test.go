package marshal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal/testutil"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

func TestTrace_Graph(t *testing.T) {
	disableColors()

	var err error
	b := &bytes.Buffer{}
	m := newTrace()

	ab := gofer.Pair{Base: "A", Quote: "B"}
	ns := testutil.Models(ab)

	err = m.Write(b, ns[ab])
	assert.NoError(t, err)

	err = m.Flush()
	assert.NoError(t, err)

	expected := `
Graph for A/B:
───median(pair:A/B)
   ├──origin(origin:a, pair:A/B)
   ├──indirect(pair:A/B)
   │  └──origin(origin:a, pair:A/B)
   └──median(pair:A/B)
      ├──origin(origin:a, pair:A/B)
      └──origin(origin:b, pair:A/B)
`[1:]

	assert.Equal(t, expected, b.String())
}

func TestTrace_Prices(t *testing.T) {
	disableColors()

	var err error
	b := &bytes.Buffer{}
	m := newTrace()

	ab := gofer.Pair{Base: "A", Quote: "B"}
	ts := testutil.Prices(ab)

	err = m.Write(b, ts[ab])
	assert.NoError(t, err)

	err = m.Flush()
	assert.NoError(t, err)

	expected := `
Price for A/B:
───aggregator(method:median, minimumSuccessfulSources:1, pair:A/B, price:10, timestamp:1970-01-01T00:00:10Z)
   ├──origin(origin:a, pair:A/B, price:10, timestamp:1970-01-01T00:00:10Z)
   ├──aggregator(method:indirect, pair:A/B, price:10, timestamp:1970-01-01T00:00:10Z)
   │  └──origin(origin:a, pair:A/B, price:10, timestamp:1970-01-01T00:00:10Z)
   └──aggregator(method:median, minimumSuccessfulSources:1, pair:A/B, price:10, timestamp:1970-01-01T00:00:10Z)
      ├──origin(origin:a, pair:A/B, price:10, timestamp:1970-01-01T00:00:10Z)
      └──origin(origin:b, pair:A/B, price:20, timestamp:1970-01-01T00:00:20Z)
            Error: something
`[1:]

	assert.Equal(t, expected, b.String())
}
