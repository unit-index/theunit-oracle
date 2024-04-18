package origins

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	pkgEthereum "github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

//go:embed curve_abi.json
var curvePoolABI string

// TODO: should be configurable
const curveDenominator = 1e18

type CurveFinance struct {
	ethClient                 pkgEthereum.Client
	addrs                     ContractAddresses
	abi                       abi.ABI
	baseIndex, quoteIndex, dx *big.Int
}

func NewCurveFinance(cli pkgEthereum.Client, addrs ContractAddresses) (*CurveFinance, error) {
	a, err := abi.JSON(strings.NewReader(curvePoolABI))
	if err != nil {
		return nil, err
	}
	return &CurveFinance{
		ethClient:  cli,
		addrs:      addrs,
		abi:        a,
		baseIndex:  big.NewInt(0),
		quoteIndex: big.NewInt(1),
		dx:         new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)),
	}, nil
}

func (s CurveFinance) pairsToContractAddress(pair Pair) (common.Address, bool, error) {
	contract, inverted, ok := s.addrs.ByPair(pair)
	if !ok {
		return common.Address{}, inverted, fmt.Errorf("failed to get contract address for pair: %s", pair.String())
	}
	return common.HexToAddress(contract), inverted, nil
}

func (s CurveFinance) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&s, pairs)
}

func (s CurveFinance) callOne(pair Pair) (*Price, error) {
	contract, inverted, err := s.pairsToContractAddress(pair)
	if err != nil {
		return nil, err
	}

	var callData []byte
	if !inverted {
		callData, err = s.abi.Pack("get_dy", s.baseIndex, s.quoteIndex, s.dx)
	} else {
		callData, err = s.abi.Pack("get_dy", s.quoteIndex, s.baseIndex, s.dx)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to pack contract args for pair: %s", pair.String())
	}

	resp, err := s.ethClient.Call(context.Background(), pkgEthereum.Call{Address: contract, Data: callData})
	if err != nil {
		return nil, err
	}
	bn := new(big.Int).SetBytes(resp)
	price, _ := new(big.Float).Quo(new(big.Float).SetInt(bn), new(big.Float).SetUint64(curveDenominator)).Float64()

	return &Price{
		Pair:      pair,
		Price:     price,
		Timestamp: time.Now(),
	}, nil
}
