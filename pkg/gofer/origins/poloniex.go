package origins

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

const poloniexURL = "https://poloniex.com/public?command=returnTicker"

type poloniexResponse struct {
	Last       stringAsFloat64 `json:"Last"`
	HidPrice   stringAsFloat64 `json:"highestBid"`
	LowestAsk  stringAsFloat64 `json:"lowestAsk"`
	BaseVolume stringAsFloat64 `json:"baseVolume"`
	IsFrozen   string          `json:"isFrozen"`
}

// Poloniex origin handler
type Poloniex struct {
	WorkerPool query.WorkerPool
}

func (p *Poloniex) localPairName(pair Pair) string {
	const (
		REP   = "REP"
		REPV2 = "REPV2"
	)

	if pair.Quote == REP {
		pair.Quote = REPV2
	}

	if pair.Base == REP {
		pair.Base = REPV2
	}

	return fmt.Sprintf("%s_%s", pair.Quote, pair.Base)
}

func (p Poloniex) Pool() query.WorkerPool {
	return p.WorkerPool
}

func (p Poloniex) PullPrices(pairs []Pair) []FetchResult {
	var err error
	req := &query.HTTPRequest{
		URL: poloniexURL,
	}

	// make query
	res := p.Pool().Query(req)
	if res == nil {
		return fetchResultListWithErrors(pairs, ErrEmptyOriginResponse)
	}
	if res.Error != nil {
		return fetchResultListWithErrors(pairs, res.Error)
	}

	// parse JSON
	var resp map[string]poloniexResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse Poloniex response: %w", err))
	}

	// prepare result
	results := make([]FetchResult, 0)
	for _, pair := range pairs {
		if r, ok := resp[p.localPairName(pair)]; !ok {
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: ErrMissingResponseForPair,
			})
		} else {
			if r.IsFrozen == "0" {
				results = append(results, FetchResult{
					Price: Price{
						Pair:      pair,
						Price:     r.Last.val(),
						Bid:       r.HidPrice.val(),
						Ask:       r.LowestAsk.val(),
						Volume24h: r.BaseVolume.val(),
						Timestamp: time.Now(),
					},
				})
			} else {
				results = append(results, FetchResult{
					Price: Price{Pair: pair},
					Error: fmt.Errorf("pair is indicated as a frozen"),
				})
			}
		}
	}

	return results
}
