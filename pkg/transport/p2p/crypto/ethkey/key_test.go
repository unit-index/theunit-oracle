package ethkey

import (
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

var (
	testAddress1 = ethereum.HexToAddress("0x2d800d93b065ce011af83f316cef9f0d005b0aa4")
	testAddress2 = ethereum.HexToAddress("0x8eb3daaf5cb4138f5f96711c09c0cfd0288a36e9")
)

func TestAddressToPeerID(t *testing.T) {
	assert.Equal(
		t,
		"1Afqz6rsuyYpr7Dpp12PbftE22nYH3k2Fw5",
		HexAddressToPeerID("0x69B352cbE6Fc5C130b6F62cc8f30b9d7B0DC27d0").Pretty(),
	)

	assert.Equal(
		t,
		"",
		HexAddressToPeerID("").Pretty(),
	)
}

func TestPeerIDToAddress(t *testing.T) {
	id, _ := peer.Decode("1Afqz6rsuyYpr7Dpp12PbftE22nYH3k2Fw5")

	assert.Equal(
		t,
		"0x69B352cbE6Fc5C130b6F62cc8f30b9d7B0DC27d0",
		PeerIDToAddress(id).String(),
	)
}
