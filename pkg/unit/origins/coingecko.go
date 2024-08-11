package origins

import (
	"encoding/json"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/internal/query"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/origins"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"strings"
)

// CoinGecko URL
const coinGeckoURL = "https://pro-api.coingecko.com/api/v3/coins/markets"

type coinGeckoResponse struct {
	CirculatingSupply string `json:"circulating_supply"`
	LastUpdated       string `json:"last_updated"`
}

// CoinGecko origin handler
type CoinGecko struct {
	WorkerPool query.WorkerPool
	Key        string
}

func (c *CoinGecko) localPairName(symbol string, name string) string {
	return fmt.Sprintf("?vs_currency=%s&ids=%s", strings.ToLower(symbol), strings.ToLower(name))
}

func (c *CoinGecko) getURL(symbol string, name string) string {
	return fmt.Sprintf(coinGeckoURL, c.localPairName(symbol, name))
}

func (c CoinGecko) Pool() query.WorkerPool {
	return c.WorkerPool
}

func (c *CoinGecko) callOne(token unit.Token) (unit.Token, error) {
	//var err error
	//req := &query.HTTPRequest{
	//	URL: c.getURL(pair),
	//}
	//
	//// make query
	//res := c.Pool().Query(req)
	//if res == nil {
	//	return &Price{}, ErrEmptyOriginResponse
	//}
	//if res.Error != nil {
	//	return &Price{}, res.Error
	//}
	//// parsing JSON
	//var resp coinGeckoResponse
	//err = json.Unmarshal(res.Body, &resp)
	//if err != nil {
	//	return &Price{}, fmt.Errorf("failed to parse CoinGecko response: %w", err)
	//}

	//getCirculatingSupply()

	//price := float64(0)
	//volume := float64(0)
	//ask := float64(0)
	//bid := float64(0)

	return token, nil
}

func (c *CoinGecko) getCirculatingSupply(symbol string, name string) (string, error) {
	var err error
	req := &query.HTTPRequest{
		URL: c.getURL(symbol, name),
	}

	// make query
	res := c.Pool().Query(req)
	if res == nil {
		return "0", origins.ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return "0", res.Error
	}
	// parsing JSON
	var resp coinGeckoResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return "0", fmt.Errorf("failed to parse CoinGecko response: %w", err)
	}

	return resp.CirculatingSupply, nil
}
