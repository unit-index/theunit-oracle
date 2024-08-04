package unit

import (
	"github.com/toknowwhy/theunit-oracle/internal/query"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/origins"
)

func NewHandler(
	origin string,
	wp query.WorkerPool,
	key string,
) (origins.Handler, error) {

	switch origin {
	case "coingecko":
		return origins.NewBaseHandler(&origins.CoinGecko{WorkerPool: wp, Key: key}), nil
	}

	return nil, nil
}
