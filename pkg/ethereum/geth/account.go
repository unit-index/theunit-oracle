package geth

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

var ErrMissingAccount = errors.New("unable to find account for requested address")

type Account struct {
	accountManager *accounts.Manager
	passphrase     string
	address        ethereum.Address
	wallet         accounts.Wallet
	account        *accounts.Account
}

// NewAccount returns a new Account instance.
func NewAccount(keyStorePath, passphrase string, address ethereum.Address) (*Account, error) {
	var err error

	if keyStorePath == "" {
		keyStorePath = defaultKeyStorePath()
	}

	ks := keystore.NewKeyStore(keyStorePath, keystore.LightScryptN, keystore.LightScryptP)

	w := &Account{
		accountManager: accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, ks),
		passphrase:     passphrase,
		address:        address,
	}

	if w.wallet, w.account, err = w.findAccountByAddress(address); err != nil {
		return nil, err
	}

	return w, nil
}

// Address returns a address of this account.
func (s *Account) Address() ethereum.Address {
	return s.address
}

// Passphrase returns a password of this account.
func (s *Account) Passphrase() string {
	return s.passphrase
}

func (s *Account) findAccountByAddress(from ethereum.Address) (accounts.Wallet, *accounts.Account, error) {
	for _, wallet := range s.accountManager.Wallets() {
		for _, account := range wallet.Accounts() {
			fmt.Println(account.Address)
			if account.Address == from {
				return wallet, &account, nil
			}
		}
	}
	//fmt.Println("aaaaaa", from)
	return nil, nil, ErrMissingAccount
}

// source: https://github.com/dapphub/dapptools/blob/master/src/ethsign/ethsign.go
func defaultKeyStorePath() string {
	var defaultKeyStores []string

	switch runtime.GOOS {
	case "darwin":
		defaultKeyStores = []string{
			os.Getenv("HOME") + "/Library/Ethereum/keystore",
			os.Getenv("HOME") + "/Library/Application Support/io.parity.ethereum/keys/ethereum",
		}
	case "windows":
		defaultKeyStores = []string{
			os.Getenv("APPDATA") + "/Ethereum/keystore",
			os.Getenv("APPDATA") + "/Parity/Ethereum/keys",
		}
	default:
		defaultKeyStores = []string{
			os.Getenv("HOME") + "/.ethereum/keystore",
			os.Getenv("HOME") + "/.local/share/io.parity.ethereum/keys/ethereum",
			os.Getenv("HOME") + "/snap/geth/current/.ethereum/keystore",
			os.Getenv("HOME") + "/snap/parity/current/.local/share/io.parity.ethereum/keys/ethereum",
		}
	}

	for _, keyStore := range defaultKeyStores {
		if _, err := os.Stat(keyStore); !os.IsNotExist(err) {
			return keyStore
		}
	}

	return ""
}
