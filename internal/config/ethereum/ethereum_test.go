package ethereum

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum/geth"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum/geth/mocks"
)

func TestEthereum_ConfigureSigner_WithoutPassword(t *testing.T) {
	config := Ethereum{
		From:     "0x07a35a1d4b751a818d93aa38e615c0df23064881",
		Keystore: "./testdata/keystore",
		Password: "",
		RPC:      "",
	}

	signer, err := config.ConfigureSigner()
	require.NoError(t, err)

	signature, err := signer.Signature([]byte("test"))
	require.NoError(t, err)

	assert.Equal(
		t,
		"b69a3cb9d029026921858b86d75f6877a0288a2b7e138076f217d3cc26e023e67e77a71ed1f5c7a1a13e0c9014d8e958d493fab36bff901b033ba1ad556df46f1c",
		hex.EncodeToString(signature.Bytes()),
	)
}

func TestEthereum_ConfigureSigner_WithPassword(t *testing.T) {
	config := Ethereum{
		From:     "2d800d93b065ce011af83f316cef9f0d005b0aa4",
		Keystore: "./testdata/keystore",
		Password: "./testdata/2.pass",
		RPC:      "",
	}

	signer, err := config.ConfigureSigner()
	require.NoError(t, err)

	signature, err := signer.Signature([]byte("test"))
	require.NoError(t, err)

	assert.Equal(
		t,
		"9c22c5f33a59a7e0d309e74ce2f448663d18d9d90b67de692a26134ba2f5cbb64826cbd408f8ce8f067205f36c614bb145c2cc1acc3902bd2d40cb1a0626a9361b",
		hex.EncodeToString(signature.Bytes()),
	)
}

func TestEthereum_ConfigureEthereumClient(t *testing.T) {
	prevEthClientFactory := ethClientFactory
	defer func() { ethClientFactory = prevEthClientFactory }()
	ethClientFactory = func(endpoints []string) (geth.EthClient, error) {
		assert.Equal(t, "1.2.3.4:1234", endpoints[0])
		return &mocks.EthClient{}, nil
	}

	config := Ethereum{
		From:     "0x07a35a1d4b751a818d93aa38e615c0df23064881",
		Keystore: "./testdata/keystore",
		Password: "",
		RPC:      "1.2.3.4:1234",
	}

	signer, err := config.ConfigureSigner()
	require.NoError(t, err)

	client, err := config.ConfigureEthereumClient(signer)
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestEthereum_ConfigureEthereumClientWithMultipleEndpoints(t *testing.T) {
	prevEthClientFactory := ethClientFactory
	defer func() { ethClientFactory = prevEthClientFactory }()
	ethClientFactory = func(endpoints []string) (geth.EthClient, error) {
		assert.Equal(t, "1.2.3.4:1234", endpoints[0])
		assert.Equal(t, "5.6.7.8:1234", endpoints[1])
		return &mocks.EthClient{}, nil
	}

	config := Ethereum{
		From:     "0x07a35a1d4b751a818d93aa38e615c0df23064881",
		Keystore: "./testdata/keystore",
		Password: "",
		RPC:      []interface{}{"1.2.3.4:1234", "5.6.7.8:1234"},
	}

	signer, err := config.ConfigureSigner()
	require.NoError(t, err)

	client, err := config.ConfigureEthereumClient(signer)
	require.NoError(t, err)
	assert.NotNil(t, client)
}
