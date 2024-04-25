package origins

import (
	"encoding/json"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/internal/query"
	"time"
)

// const okexURL = "https://www.okex.com/api/spot/v3/instruments/ticker"
const okexURL = "https://www.okx.com/api/v5/market/tickers?instType=SPOT"

type okexResponse struct {
	InstrumentID  string          `json:"instId"`
	Last          stringAsFloat64 `json:"last"`
	BestAsk       stringAsFloat64 `json:"askPx"`
	BestBid       stringAsFloat64 `json:"bidPx"`
	BaseVolume24H stringAsFloat64 `json:"volCcy24h"`
	Timestamp     string          `json:"ts"`
}

type okexV5Response struct {
	Code string         `json:"code"`
	Data []okexResponse `json:"data"`
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
	//var resp []okexResponse
	resp := okexV5Response{}
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse Okex response: %w", err))
	}

	// convert response from a slice to a map
	respMap := map[string]okexResponse{}

	for _, symbolResp := range resp.Data {
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
			//ti, err := strconv.Atoi(r.Timestamp)
			//if err != nil {
			//	//fmt.Println("Error converting string to int:", err)
			//	return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse Okex response: %w", err))
			//}
			//t := time.Unix(int64(ti), 0)

			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     r.Last.val(),
					Bid:       r.BestBid.val(),
					Ask:       r.BestAsk.val(),
					Volume24h: r.BaseVolume24H.val(),
					Timestamp: time.Now(),
				},
			})
		}
	}

	return results
}
