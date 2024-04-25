package feeds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

func TestFeeds_Addresses_Valid(t *testing.T) {
	feeds := Feeds{"0x07a35a1d4b751a818d93aa38e615c0df23064881", "2d800d93b065ce011af83f316cef9f0d005b0aa4"}
	addrs, err := feeds.Addresses()
	require.NoError(t, err)

	assert.Equal(t, ethereum.HexToAddress("0x07a35a1d4b751a818d93aa38e615c0df23064881"), addrs[0])
	assert.Equal(t, ethereum.HexToAddress("0x2d800d93b065ce011af83f316cef9f0d005b0aa4"), addrs[1])
}

func TestFeeds_Addresses_Invalid(t *testing.T) {
	feeds := Feeds{"0x07a35a1d4b751a818d93aa38e615c0df23064881", "abc"}
	_, err := feeds.Addresses()

	require.ErrorIs(t, err, ErrInvalidEthereumAddress)
}
