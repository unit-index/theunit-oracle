package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTypeValue(t *testing.T) {
	for ct, st := range formatMap {
		t.Run(st, func(t *testing.T) {
			ftv := formatTypeValue{}
			err := ftv.Set(st)

			assert.Nil(t, err)
			assert.Equal(t, st, ftv.String())
			assert.Equal(t, ct, ftv.format)
		})
	}
}
