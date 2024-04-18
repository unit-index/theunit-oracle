package origins

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

//go:embed wsteth_abi.json
var wrappedStakedETHABI string

const wsethDenominator = 1e18

type WrappedStakedETH struct {
	ethClient ethereum.Client
	addrs     ContractAddresses
	abi       abi.ABI
}

func NewWrappedStakedETH(cli ethereum.Client, addrs ContractAddresses) (*WrappedStakedETH, error) {
	a, err := abi.JSON(strings.NewReader(wrappedStakedETHABI))
	if err != nil {
		return nil, err
	}
	return &WrappedStakedETH{
		ethClient: cli,
		addrs:     addrs,
		abi:       a,
	}, nil
}

func (s WrappedStakedETH) pairsToContractAddress(pair Pair) (ethereum.Address, bool, error) {
	contract, inverted, ok := s.addrs.ByPair(pair)
	if !ok {
		return ethereum.Address{}, inverted, fmt.Errorf("failed to get Curve contract address for pair: %s", pair.String())
	}
	return ethereum.HexToAddress(contract), inverted, nil
}

func (s WrappedStakedETH) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&s, pairs)
}

func (s WrappedStakedETH) callOne(pair Pair) (*Price, error) {
	contract, inverted, err := s.pairsToContractAddress(pair)
	if err != nil {
		return nil, err
	}

	var callData []byte
	if !inverted {
		callData, err = s.abi.Pack("stEthPerToken")
	} else {
		callData, err = s.abi.Pack("tokensPerStEth")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get contract args for pair: %s", pair.String())
	}

	resp, err := s.ethClient.Call(context.Background(), ethereum.Call{Address: contract, Data: callData})
	if err != nil {
		return nil, err
	}
	bn := new(big.Int).SetBytes(resp)
	price, _ := new(big.Float).Quo(new(big.Float).SetInt(bn), new(big.Float).SetUint64(wsethDenominator)).Float64()

	return &Price{
		Pair:      pair,
		Price:     price,
		Timestamp: time.Now(),
	}, nil
}
