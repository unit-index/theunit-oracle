package origins

import (
	"encoding/json"
	"fmt"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

type Kyber struct {
	WorkerPool query.WorkerPool
}

func (k Kyber) Pool() query.WorkerPool {
	return k.WorkerPool
}

func (k Kyber) PullPrices(pairs []Pair) []FetchResult {
	req := &query.HTTPRequest{
		URL: kyberURL,
	}
	res := k.Pool().Query(req)
	if errorResponses := validateResponse(pairs, res); len(errorResponses) > 0 {
		return errorResponses
	}
	return k.parseResponse(pairs, res)
}

const kyberURL = "https://api.kyber.network/change24h"

type kyberTicker struct {
	Timestamp    intAsUnixTimestampMs `json:"timestamp"`
	TokenName    string               `json:"token_name"`
	TokenSymbol  string               `json:"token_symbol"`
	TokenDecimal int                  `json:"token_decimal"`
	TokenAddress string               `json:"token_address"`
	RateEthNow   float64              `json:"rate_eth_now"`
	ChangeEth24H float64              `json:"change_eth_24h"`
	ChangeUsd24H float64              `json:"change_usd_24h"`
	RateUsdNow   float64              `json:"rate_usd_now"`
}

func (k *Kyber) parseResponse(pairs []Pair, res *query.HTTPResponse) []FetchResult {
	results := make([]FetchResult, 0)
	var tickers map[string]kyberTicker
	err := json.Unmarshal(res.Body, &tickers)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse response: %w", err))
	}

	for _, pair := range pairs {
		//nolint:gocritic
		if t, is := tickers[pair.Quote+"_"+pair.Base]; !is {
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: ErrMissingResponseForPair,
			})
		} else if t.TokenSymbol != pair.Base {
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: ErrInvalidPrice,
			})
		} else {
			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     t.RateEthNow,
					Timestamp: t.Timestamp.val(),
				},
			})
		}
	}
	return results
}
