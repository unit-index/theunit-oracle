package origins

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

type Kraken struct {
	WorkerPool query.WorkerPool
}

const krakenURL = "https://api.kraken.com/0/public/Ticker?pair=%s"

func (k Kraken) Pool() query.WorkerPool {
	return k.WorkerPool
}
func (k Kraken) PullPrices(pairs []Pair) []FetchResult {
	req := &query.HTTPRequest{
		URL: fmt.Sprintf(krakenURL, k.localPairName(pairs...)),
	}
	res := k.Pool().Query(req)
	if errorResponses := validateResponse(pairs, res); len(errorResponses) > 0 {
		return errorResponses
	}
	return k.parseResponse(pairs, res)
}

type krakenResponse struct {
	Errors []string `json:"error"`
	Result map[string]krakenPairResponse
}

type krakenPairResponse struct {
	Price  firstStringFromSliceAsFloat64 `json:"c"`
	Volume firstStringFromSliceAsFloat64 `json:"v"`
	Ask    firstStringFromSliceAsFloat64 `json:"a"`
	Bid    firstStringFromSliceAsFloat64 `json:"b"`
}

func (k *Kraken) parseResponse(pairs []Pair, res *query.HTTPResponse) []FetchResult {
	var resp krakenResponse
	err := json.Unmarshal(res.Body, &resp)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse response: %w", err))
	}
	results := make([]FetchResult, 0)
	for _, pair := range pairs {
		if t, is := resp.Result[k.localPairName(pair)]; !is {
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: ErrMissingResponseForPair,
			})
		} else {
			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     t.Price.val(),
					Ask:       t.Ask.val(),
					Bid:       t.Bid.val(),
					Volume24h: t.Volume.val(),
					Timestamp: time.Now(),
				},
			})
		}
	}
	return results
}

func (k *Kraken) localPairName(pairs ...Pair) string {
	var l []string
	for _, pair := range pairs {
		l = append(l, pair.String())
	}
	return strings.Join(l, ",")
}
