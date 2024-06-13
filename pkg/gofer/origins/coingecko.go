// TODO: check if it's possible to merge coinbase and coinbasepro
package origins

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

// Coinbase URL
const coinGeckoURL = "https://pro-api.coingecko.com/api/v3/coins/markets?vs_currency=btc&ids=bitcoin&per_page=1"

type coinGeckoResponse struct {
	CirculatingSupply string `json:"circulating_supply"`
	LastUpdated       string `json:"last_updated"`
}

// Coinbase origin handler
type CoinGecko struct {
	WorkerPool query.WorkerPool
}

func (c *CoinGecko) localPairName(pair Pair) string {
	return fmt.Sprintf("%s-%s", strings.ToUpper(pair.Base), strings.ToUpper(pair.Quote))
}

func (c *CoinGecko) getURL(pair Pair) string {
	return fmt.Sprintf(coinbaseProURL, c.localPairName(pair))
}

func (c CoinGecko) Pool() query.WorkerPool {
	return c.WorkerPool
}

func (c CoinGecko) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&c, pairs)
}

func (c *CoinGecko) callOne(pair Pair) (*Price, error) {
	var err error
	req := &query.HTTPRequest{
		URL: c.getURL(pair),
	}

	// make query
	res := c.Pool().Query(req)
	if res == nil {
		return nil, ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return nil, res.Error
	}
	// parsing JSON
	var resp coinGeckoResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coinbasepro response: %w", err)
	}

	return &Price{
		Pair:      pair,
		Timestamp: time.Now(),
	}, nil
}
