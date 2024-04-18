package marshal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal/testutil"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

func TestPlain_Nodes(t *testing.T) {
	var err error
	b := &bytes.Buffer{}
	m := newPlain()

	ab := gofer.Pair{Base: "A", Quote: "B"}
	cd := gofer.Pair{Base: "C", Quote: "D"}
	ns := testutil.Models(ab, cd)

	err = m.Write(b, ns[ab])
	assert.NoError(t, err)

	err = m.Write(b, ns[cd])
	assert.NoError(t, err)

	err = m.Flush()
	assert.NoError(t, err)

	expected := `
A/B
C/D
`[1:]

	assert.Equal(t, expected, b.String())
}

func TestPlain_Prices(t *testing.T) {
	var err error
	b := &bytes.Buffer{}
	m := newPlain()

	ab := gofer.Pair{Base: "A", Quote: "B"}
	cd := gofer.Pair{Base: "C", Quote: "D"}
	ns := testutil.Prices(ab, cd)

	err = m.Write(b, ns[ab])
	assert.NoError(t, err)

	cdt := ns[cd]
	cdt.Error = "something"
	err = m.Write(b, ns[cd])
	assert.NoError(t, err)

	err = m.Flush()
	assert.NoError(t, err)

	expected := `
A/B 10.000000
C/D - something
`[1:]

	assert.Equal(t, expected, b.String())
}
