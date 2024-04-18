package origins

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

const bittrexURL = "https://api.bittrex.com/api/v1.1/public/getticker?market=%s"

type bittrexResponse struct {
	Success bool                  `json:"success"`
	Result  bittrexSymbolResponse `json:"result"`
}

type bittrexSymbolResponse struct {
	Ask  float64 `json:"Ask"`
	Bid  float64 `json:"Bid"`
	Last float64 `json:"Last"`
}

// Bittrex origin handler
type Bittrex struct {
	WorkerPool query.WorkerPool
}

func (b *Bittrex) localPairName(pair Pair) string {
	return fmt.Sprintf("%s-%s", pair.Quote, pair.Base)
}

func (b Bittrex) Pool() query.WorkerPool {
	return b.WorkerPool
}

func (b Bittrex) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&b, pairs)
}

func (b *Bittrex) callOne(pair Pair) (*Price, error) {
	var err error
	req := &query.HTTPRequest{
		URL: fmt.Sprintf(bittrexURL, b.localPairName(pair)),
	}

	// make query
	res := b.Pool().Query(req)
	if res == nil {
		return nil, ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return nil, res.Error
	}

	// parse JSON
	var resp bittrexResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Bittrex response: %w", err)
	}
	if !resp.Success {
		return nil, fmt.Errorf("wrong response from Bittrex %v", resp)
	}

	return &Price{
		Pair:      pair,
		Price:     resp.Result.Last,
		Bid:       resp.Result.Bid,
		Ask:       resp.Result.Ask,
		Timestamp: time.Now(),
	}, nil
}
