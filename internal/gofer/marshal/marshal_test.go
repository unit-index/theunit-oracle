package marshal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMarshaller(t *testing.T) {
	expectedMap := map[FormatType]interface{}{
		Plain:  (*plain)(nil),
		JSON:   (*json)(nil),
		NDJSON: (*json)(nil),
		Trace:  (*trace)(nil),
	}
	formatMap := map[FormatType]string{
		Plain:  "plain",
		JSON:   "json",
		NDJSON: "ndjson",
		Trace:  "trace",
	}
	for ct, st := range formatMap {
		t.Run(st, func(t *testing.T) {
			m, err := NewMarshal(ct)

			assert.NoError(t, err)
			assert.Implements(t, (*Marshaller)(nil), m)
			assert.IsType(t, m, &Marshal{})
			assert.IsType(t, m.marshaller, expectedMap[ct])
		})
	}
}
