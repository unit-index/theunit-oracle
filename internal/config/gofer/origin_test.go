package gofer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingOriginParamsAliasesFailParsing(t *testing.T) {
	parsed, err := parseParamsSymbolAliases(nil)
	assert.Nil(t, parsed)
	assert.Error(t, err)

	parsed, err = parseParamsSymbolAliases([]byte(""))
	assert.Nil(t, parsed)
	assert.Error(t, err)
}

func TestParsingOriginParamsAliases(t *testing.T) {
	// parsing empty aliases
	parsed, err := parseParamsSymbolAliases([]byte(`{}`))
	assert.NoError(t, err)
	assert.Nil(t, parsed)

	// Parsing only apiKey
	key, err := parseParamsAPIKey([]byte(`{"apiKey":"test"}`))
	assert.NoError(t, err)
	assert.Equal(t, "test", key)

	// Parsing contracts
	contracts, err := parseParamsContracts([]byte(`{"contracts":{"BTC/ETH":"0x00000"}}`))
	assert.NoError(t, err)
	assert.NotNil(t, contracts)
	assert.Equal(t, "0x00000", contracts["BTC/ETH"])

	// Parsing symbol aliases
	aliases, err := parseParamsSymbolAliases([]byte(`{"symbolAliases":{"ETH":"WETH"}}`))
	assert.NoError(t, err)
	assert.NotNil(t, aliases)
	assert.Equal(t, "WETH", aliases["ETH"])
}
