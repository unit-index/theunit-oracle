package origins

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

// Fx URL
const fxURL = "https://api.exchangeratesapi.io/latest?symbols=%s&base=%s&access_key=%s"

type fxResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// Fx exchange handler
type Fx struct {
	WorkerPool query.WorkerPool
	APIKey     string
}

func (f *Fx) renameSymbol(symbol string) string {
	return strings.ToUpper(symbol)
}

func (f Fx) Pool() query.WorkerPool {
	return f.WorkerPool
}

func (f Fx) PullPrices(pairs []Pair) []FetchResult {
	// Group pairs by asset pair base.
	bases := map[string][]Pair{}
	for _, pair := range pairs {
		base := pair.Base
		bases[base] = append(bases[base], pair)
	}

	var results []FetchResult
	for base, pairs := range bases {
		// Make one request per asset pair base.
		crs, err := f.callByBase(base, pairs)
		if err != nil {
			// If callByBase fails wholesale, create a FetchResult per pair with the same
			// error.
			crs = fetchResultListWithErrors(pairs, err)
		}
		results = append(results, crs...)
	}

	return results
}

func (f *Fx) getURL(base string, quotes []Pair) string {
	symbols := []string{}
	for _, pair := range quotes {
		symbols = append(symbols, f.renameSymbol(pair.Quote))
	}
	return fmt.Sprintf(fxURL, strings.Join(symbols, ","), f.renameSymbol(base), f.APIKey)
}

func (f *Fx) callByBase(base string, pairs []Pair) ([]FetchResult, error) {
	req := &query.HTTPRequest{
		URL: f.getURL(base, pairs),
	}

	// Make query.
	res := f.Pool().Query(req)
	if res == nil {
		return nil, ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return nil, res.Error
	}
	// Parse JSON.
	var resp fxResponse
	err := json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse FX response: %w", err)
	}
	if resp.Rates == nil {
		return nil, fmt.Errorf("failed to parse FX response: %+v", resp)
	}

	results := make([]FetchResult, len(pairs))
	for i, pair := range pairs {
		if price, ok := resp.Rates[f.renameSymbol(pair.Quote)]; ok {
			// Build Price from exchange response.
			results[i] = FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     price,
					Timestamp: time.Now(),
				},
				Error: nil,
			}
		} else {
			// Missing quote in exchange response.
			results[i] = fetchResultWithError(
				pair,
				fmt.Errorf("no price for %s quote exist in response %s", pair.Quote, res.Body),
			)
		}
	}
	return results, nil
}
