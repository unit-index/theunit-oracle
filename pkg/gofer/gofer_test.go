package gofer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPair(t *testing.T) {
	tests := []struct {
		name    string
		pair    string
		want    Pair
		wantErr bool
	}{
		{
			name:    "valid-pair",
			pair:    "A/B",
			want:    Pair{Base: "A", Quote: "B"},
			wantErr: false,
		},
		{
			name:    "valid-lowercase-pair",
			pair:    "a/b",
			want:    Pair{Base: "A", Quote: "B"},
			wantErr: false,
		},
		{
			name:    "missing-slash",
			pair:    "AB",
			want:    Pair{},
			wantErr: false,
		},
		{
			name:    "multiple-slashes",
			pair:    "A/B/",
			want:    Pair{},
			wantErr: false,
		},
		{
			name:    "empty",
			pair:    "",
			want:    Pair{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPair(tt.pair)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestNewPairs(t *testing.T) {
	tests := []struct {
		name    string
		pairs   []string
		want    []Pair
		wantErr bool
	}{
		{
			name:    "single-valid-pair",
			pairs:   []string{"A/B"},
			want:    []Pair{{Base: "A", Quote: "B"}},
			wantErr: false,
		},
		{
			name:    "multiple-valid-pair",
			pairs:   []string{"A/B", "X/Y"},
			want:    []Pair{{Base: "A", Quote: "B"}, {Base: "X", Quote: "Y"}},
			wantErr: false,
		},
		{
			name:    "contains-invalid-pair",
			pairs:   []string{"A/B", "XY"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			pairs:   []string{},
			want:    []Pair(nil),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPairs(tt.pairs...)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
