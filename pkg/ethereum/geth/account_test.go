package geth

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var accountAddress = common.HexToAddress("0x2d800d93b065ce011af83f316cef9f0d005b0aa4")

func TestAccount_ValidAddress(t *testing.T) {
	account, err := NewAccount("./testdata/keystore", "test123", accountAddress)
	assert.NoError(t, err)

	assert.Equal(t, accountAddress, account.Address())
	assert.Equal(t, "test123", account.Passphrase())
	assert.Equal(t, accountAddress, account.account.Address)
	assert.NotNil(t, account.wallet)
}

func TestAccount_InvalidAddress(t *testing.T) {
	account, err := NewAccount("./testdata/keystore", "test123", common.HexToAddress(""))
	assert.Error(t, err)
	assert.Nil(t, account)
}
