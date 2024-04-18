package origins

import (
	"encoding/json"
	"fmt"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

const binanceURL = "https://www.binance.com/api/v3/ticker/24hr"

type binanceResponse struct {
	Symbol    string               `json:"symbol"`
	LastPrice stringAsFloat64      `json:"lastPrice"`
	BidPrice  stringAsFloat64      `json:"bidPrice"`
	AskPrice  stringAsFloat64      `json:"askPrice"`
	Volume    stringAsFloat64      `json:"volume"`
	CloseTime intAsUnixTimestampMs `json:"closeTime"`
}

// Binance origin handler
type Binance struct {
	WorkerPool query.WorkerPool
}

func (b *Binance) localPairName(pair Pair) string {
	return pair.Base + pair.Quote
}

func (b Binance) Pool() query.WorkerPool {
	return b.WorkerPool
}

func (b Binance) PullPrices(pairs []Pair) []FetchResult {
	var err error
	req := &query.HTTPRequest{
		URL: binanceURL,
	}

	// make query
	res := b.WorkerPool.Query(req)
	if res == nil {
		return fetchResultListWithErrors(pairs, ErrEmptyOriginResponse)
	}
	if res.Error != nil {
		return fetchResultListWithErrors(pairs, res.Error)
	}

	// parse JSON
	var resp []binanceResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse Binance response: %w", err))
	}

	// convert response from a slice to a map
	respMap := map[string]binanceResponse{}
	for _, symbolResp := range resp {
		respMap[symbolResp.Symbol] = symbolResp
	}

	// prepare result
	results := make([]FetchResult, 0)
	for _, pair := range pairs {
		if r, ok := respMap[b.localPairName(pair)]; !ok {
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: ErrMissingResponseForPair,
			})
		} else {
			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     r.LastPrice.val(),
					Bid:       r.BidPrice.val(),
					Ask:       r.AskPrice.val(),
					Volume24h: r.Volume.val(),
					Timestamp: r.CloseTime.val(),
				},
			})
		}
	}

	return results
}
