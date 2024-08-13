package geth

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"math/big"
)

type UnitAlgorithm struct {
	ethereum ethereum.Client
	address  ethereum.Address
}
type MarketInfo struct {
	lastMonthClosePrice float64
	lastMonthMarketCap  float64
}

// NewMedian creates the new Median instance.
func NewUnitAlgorithm(ethereum ethereum.Client, address ethereum.Address) *UnitAlgorithm {
	return &UnitAlgorithm{
		ethereum: ethereum,
		address:  address,
	}
}
func (u *UnitAlgorithm) GetTokens(ctx context.Context, time *big.Int) ([]common.Address, error) {
	r, err := u.read(ctx, "getTokens", time)
	if err != nil {
		return nil, err
	}
	b := r[0].([]common.Address)
	return b, nil
}

func (u *UnitAlgorithm) TokenPerMonthMarketInfo(ctx context.Context) ([]string, error) {
	return nil, nil

}
func (u *UnitAlgorithm) TotalMarketCap(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (u *UnitAlgorithm) read(ctx context.Context, method string, args ...interface{}) ([]interface{}, error) {
	cd, err := unitAlgorithmABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}

	var data []byte
	err = retry(maxReadRetries, delayBetweenReadRetries, func() error {
		data, err = u.ethereum.Call(ctx, ethereum.Call{Address: u.address, Data: cd})
		return err
	})
	if err != nil {
		return nil, err
	}

	return unitAlgorithmABI.Unpack(method, data)
}
