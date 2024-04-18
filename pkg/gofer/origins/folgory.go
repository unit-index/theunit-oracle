package origins

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

type Folgory struct {
	WorkerPool query.WorkerPool
}

func (o Folgory) Pool() query.WorkerPool {
	return o.WorkerPool
}

func (o Folgory) PullPrices(pairs []Pair) []FetchResult {
	req := &query.HTTPRequest{
		URL: folgoryURL,
	}
	res := o.Pool().Query(req)
	if errorResponses := validateResponse(pairs, res); len(errorResponses) > 0 {
		return errorResponses
	}
	return o.parseResponse(pairs, res)
}

const folgoryURL = "https://folgory.com/api/v1"

type folgoryTicker struct {
	Symbol string          `json:"symbol"`
	Price  stringAsFloat64 `json:"last"`
	Volume stringAsFloat64 `json:"volume"`
}

func (o *Folgory) localPairName(pair Pair) string {
	return fmt.Sprintf("%s/%s", pair.Base, pair.Quote)
}

func (o *Folgory) parseResponse(pairs []Pair, res *query.HTTPResponse) []FetchResult {
	results := make([]FetchResult, 0)
	var resp []folgoryTicker
	err := json.Unmarshal(res.Body, &resp)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse response: %w", err))
	}

	tickers := make(map[string]folgoryTicker)
	for _, t := range resp {
		tickers[t.Symbol] = t
	}

	for _, pair := range pairs {
		if t, is := tickers[o.localPairName(pair)]; !is {
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: fmt.Errorf("no response for %s", pair.String()),
			})
		} else {
			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     t.Price.val(),
					Volume24h: t.Volume.val(),
					Timestamp: time.Now(),
				},
			})
		}
	}
	return results
}
