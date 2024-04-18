package origins

import (
	"encoding/json"
	"fmt"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

// Coinbase URL
const openExchangeRatesURL = "https://openexchangerates.org/api/latest.json?app_id=%s&base=%s&symbols=%s"

type openExchangeRatesResponse struct {
	Timestamp intAsUnixTimestamp `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}

// OpenExchangeRates origin handler
type OpenExchangeRates struct {
	WorkerPool query.WorkerPool
	APIKey     string
}

func (o *OpenExchangeRates) getURL(pair Pair) string {
	return fmt.Sprintf(openExchangeRatesURL, o.APIKey, pair.Base, pair.Quote)
}

func (o OpenExchangeRates) Pool() query.WorkerPool {
	return o.WorkerPool
}
func (o OpenExchangeRates) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&o, pairs)
}

func (o *OpenExchangeRates) callOne(pair Pair) (*Price, error) {
	var err error
	req := &query.HTTPRequest{
		URL: o.getURL(pair),
	}

	// make query
	res := o.Pool().Query(req)
	if res == nil {
		return nil, ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return nil, res.Error
	}
	// parsing JSON
	var resp openExchangeRatesResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenExchangeRate response: %w", err)
	}
	price, ok := resp.Rates[pair.Quote]
	if !ok {
		return nil, ErrMissingResponseForPair
	}
	// building Price
	return &Price{
		Pair:      pair,
		Price:     price,
		Timestamp: resp.Timestamp.val(),
	}, nil
}
