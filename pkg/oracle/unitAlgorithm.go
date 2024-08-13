package oracle

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type UnitAlgorithm interface {
	GetTokens(ctx context.Context, time *big.Int) ([]common.Address, error)
}
