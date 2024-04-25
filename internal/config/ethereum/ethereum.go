package ethereum

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/toknowwhy/theunit-oracle/internal/rpcsplitter"
	"github.com/toknowwhy/theunit-oracle/pkg/log/null"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum/geth"
)

const splitterVirtualHost = "unit-splitter"

var ethClientFactory = func(endpoints []string) (geth.EthClient, error) {
	switch len(endpoints) {
	case 0:
		return nil, errors.New("missing address to a RPC client in the configuration file")
	case 1:
		return ethclient.Dial(endpoints[0])
	default:
		// TODO: pass logger
		splitter, err := rpcsplitter.NewTransport(endpoints, splitterVirtualHost, nil, null.New())

		if err != nil {
			return nil, err
		}
		rpcClient, err := rpc.DialHTTPWithClient(
			fmt.Sprintf("http://%s", splitterVirtualHost),
			&http.Client{Transport: splitter},
		)

		if err != nil {
			return nil, err
		}
		return ethclient.NewClient(rpcClient), nil
	}

}

type Ethereum struct {
	From     string      `json:"from"`
	Keystore string      `json:"keystore"`
	Password string      `json:"password"`
	RPC      interface{} `json:"rpc"`
}

func (c *Ethereum) ConfigureSigner() (ethereum.Signer, error) {
	account, err := c.configureAccount()
	if err != nil {
		return nil, err
	}
	return geth.NewSigner(account), nil
}

func (c *Ethereum) ConfigureEthereumClient(signer ethereum.Signer) (*geth.Client, error) {
	var endpoints []string
	switch v := c.RPC.(type) {
	case string:
		endpoints = []string{v}
	case []interface{}:
		for _, s := range v {
			if s, ok := s.(string); ok {
				endpoints = append(endpoints, s)
			}
		}
	}
	if len(endpoints) == 0 {
		return nil, errors.New("value of the RPC key must be string or array of strings")
	}
	client, err := ethClientFactory(endpoints)
	if err != nil {
		return nil, err
	}
	return geth.NewClient(client, signer), nil
}

func (c *Ethereum) configureAccount() (*geth.Account, error) {
	if c.From == "" {
		return nil, nil
	}
	fmt.Println(c.Password)
	//passphrase, err := c.readAccountPassphrase(c.Password)
	//if err != nil {
	//	return nil, err
	//}
	passphrase := c.Password

	account, err := geth.NewAccount(c.Keystore, passphrase, ethereum.HexToAddress(c.From))
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (c *Ethereum) readAccountPassphrase(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	passphrase, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read Ethereum password file: %w", err)
	}
	return strings.TrimSuffix(string(passphrase), "\n"), nil
}
