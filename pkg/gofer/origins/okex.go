package origins

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

const okexURL = "https://www.okex.com/api/spot/v3/instruments/ticker"

type okexResponse struct {
	InstrumentID  string          `json:"instrument_id"`
	Last          stringAsFloat64 `json:"last"`
	BestAsk       stringAsFloat64 `json:"best_ask"`
	BestBid       stringAsFloat64 `json:"best_bid"`
	BaseVolume24H stringAsFloat64 `json:"base_volume_24h"`
	Timestamp     time.Time       `json:"timestamp"`
}

// Okex origin handler
type Okex struct {
	WorkerPool query.WorkerPool
}

func (o *Okex) localPairName(pair Pair) string {
	return fmt.Sprintf("%s-%s", pair.Base, pair.Quote)
}

func (o Okex) Pool() query.WorkerPool {
	return o.WorkerPool
}
func (o Okex) PullPrices(pairs []Pair) []FetchResult {
	var err error
	req := &query.HTTPRequest{
		URL: okexURL,
	}

	// make query
	res := o.Pool().Query(req)
	if res == nil {
		return fetchResultListWithErrors(pairs, ErrEmptyOriginResponse)
	}
	if res.Error != nil {
		return fetchResultListWithErrors(pairs, res.Error)
	}

	// parse JSON
	var resp []okexResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse Okex response: %w", err))
	}

	// convert response from a slice to a map
	respMap := map[string]okexResponse{}
	for _, symbolResp := range resp {
		respMap[symbolResp.InstrumentID] = symbolResp
	}

	// prepare result
	results := make([]FetchResult, 0)
	for _, pair := range pairs {
		if r, ok := respMap[o.localPairName(pair)]; !ok {
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: ErrMissingResponseForPair,
			})
		} else {
			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     r.Last.val(),
					Bid:       r.BestBid.val(),
					Ask:       r.BestAsk.val(),
					Volume24h: r.BaseVolume24H.val(),
					Timestamp: r.Timestamp,
				},
			})
		}
	}

	return results
}
