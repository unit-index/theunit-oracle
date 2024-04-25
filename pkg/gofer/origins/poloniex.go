package origins

import (
	"encoding/json"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/internal/query"
	"strconv"
	"time"
)

// const poloniexURL = "https://poloniex.com/public?command=returnTicker"
const poloniexURL = "https://api.poloniex.com/markets/ticker24h"

type poloniexResponse struct {
	MarkPrice string             `json:"markPrice"`
	Bid       stringAsFloat64    `json:"bid"`
	Ask       stringAsFloat64    `json:"ask"`
	Quantity  stringAsFloat64    `json:"quantity"`
	Symbol    string             `json:"symbol"`
	TimeStamp intAsUnixTimestamp `json:"ts"`
}

// Poloniex origin handler
type Poloniex struct {
	WorkerPool query.WorkerPool
}

func (p *Poloniex) localPairName(pair Pair) string {
	return fmt.Sprintf("%s_%s", pair.Base, pair.Quote)
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
	var resp []poloniexResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse Poloniex response: %w", err))
	}
	// prepare result
	results := make([]FetchResult, 0)
	//for _, pair := range pairs {
	//	if r, ok := resp[p.localPairName(pair)]; !ok {
	//		results = append(results, FetchResult{
	//			Price: Price{Pair: pair},
	//			Error: ErrMissingResponseForPair,
	//		})
	//	} else {
	//		if r.IsFrozen == "0" {
	//			results = append(results, FetchResult{
	//				Price: Price{
	//					Pair:      pair,
	//					Price:     r.Last.val(),
	//					Bid:       r.HidPrice.val(),
	//					Ask:       r.LowestAsk.val(),
	//					Volume24h: r.BaseVolume.val(),
	//					Timestamp: time.Now(),
	//				},
	//			})
	//		} else {
	//			results = append(results, FetchResult{
	//				Price: Price{Pair: pair},
	//				Error: fmt.Errorf("pair is indicated as a frozen"),
	//			})
	//		}
	//	}
	//}

	respMap := map[string]poloniexResponse{}
	for _, symbolResp := range resp {
		respMap[symbolResp.Symbol] = symbolResp
	}

	for _, pair := range pairs {
		if r, ok := respMap[p.localPairName(pair)]; !ok {
			//fmt.Println("BTC***********", pair, p.localPairName(pair))
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: ErrMissingResponseForPair,
			})
		} else {
			if r.MarkPrice == "" {
				//fmt.Println("BTC***********", pair)
				results = append(results, FetchResult{
					Price: Price{Pair: pair},
					Error: ErrMissingResponseForPair,
				})
				continue
			}

			markPrice, err := strconv.ParseFloat(r.MarkPrice, 64)

			if err != nil {
				results = append(results, FetchResult{
					Price: Price{Pair: pair},
					Error: ErrMissingResponseForPair,
				})
				continue
			}
			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     markPrice,
					Bid:       r.Bid.val(),
					Ask:       r.Ask.val(),
					Volume24h: r.Quantity.val(),
					Timestamp: time.Now(),
				},
			})
		}
	}

	return results
}
