// TODO: check if it's possible to merge coinbase and coinbasepro
package origins

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

// Coinbase URL
const coinbaseProURL = "https://api.pro.coinbase.com/products/%s/ticker"

type coinbaseProResponse struct {
	Price  string `json:"price"`
	Ask    string `json:"ask"`
	Bid    string `json:"bid"`
	Volume string `json:"volume"`
}

// Coinbase origin handler
type CoinbasePro struct {
	WorkerPool query.WorkerPool
}

func (c *CoinbasePro) localPairName(pair Pair) string {
	return fmt.Sprintf("%s-%s", strings.ToUpper(pair.Base), strings.ToUpper(pair.Quote))
}

func (c *CoinbasePro) getURL(pair Pair) string {
	return fmt.Sprintf(coinbaseProURL, c.localPairName(pair))
}

func (c CoinbasePro) Pool() query.WorkerPool {
	return c.WorkerPool
}

func (c CoinbasePro) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&c, pairs)
}

func (c *CoinbasePro) callOne(pair Pair) (*Price, error) {
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
	var resp coinbaseProResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coinbasepro response: %w", err)
	}
	// Parsing price from string
	price, err := strconv.ParseFloat(resp.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price from coinbasepro origin %s", res.Body)
	}
	// Parsing ask from string
	ask, err := strconv.ParseFloat(resp.Ask, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ask from coinbasepro origin %s", res.Body)
	}
	// Parsing volume from string
	volume, err := strconv.ParseFloat(resp.Volume, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse volume from coinbasepro origin %s", res.Body)
	}
	// Parsing bid from string
	bid, err := strconv.ParseFloat(resp.Bid, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bid from coinbasepro origin %s", res.Body)
	}
	// building Price
	return &Price{
		Pair:      pair,
		Price:     price,
		Volume24h: volume,
		Ask:       ask,
		Bid:       bid,
		Timestamp: time.Now(),
	}, nil
}
