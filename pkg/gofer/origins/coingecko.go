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
const coingeckoProURL = "https://api.pro.coinbase.com/products/%s/ticker"

type coinbase2ProResponse struct {
	Price  string `json:"price"`
	Ask    string `json:"ask"`
	Bid    string `json:"bid"`
	Volume string `json:"volume"`
}

// Coinbase origin handler
type Coinbase2Pro struct {
	WorkerPool query.WorkerPool
}

func (c *Coinbase2Pro) localPairName(pair Pair) string {
	return fmt.Sprintf("%s-%s", strings.ToUpper(pair.Base), strings.ToUpper("USD"))
}

func (c *Coinbase2Pro) getURL(pair Pair) string {
	return fmt.Sprintf(coingeckoProURL, c.localPairName(pair))
}

func (c Coinbase2Pro) Pool() query.WorkerPool {
	return c.WorkerPool
}

func (c Coinbase2Pro) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&c, pairs)
}

func (c *Coinbase2Pro) callOne(pair Pair) (*Price, error) {
	fmt.Println("callOne", pair)
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
	var resp coinbase2ProResponse
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
