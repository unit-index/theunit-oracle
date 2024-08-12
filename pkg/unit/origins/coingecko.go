package origins

import (
	"encoding/json"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/internal/query"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/origins"
	"strings"
	"time"
)

// CoinGecko URL
const coinGeckoURL = "https://pro-api.coingecko.com/api/v3/coins/markets"

type coinGeckoResponse struct {
	CirculatingSupply float64 `json:"circulating_supply"`
	LastUpdated       string  `json:"last_updated"`
}

// CoinGecko origin handler
type CoinGecko struct {
	WorkerPool query.WorkerPool
	Key        string
}

func (c *CoinGecko) localPairName(symbol string, name string) string {
	return fmt.Sprintf("vs_currency=%s&ids=%s", symbol, name)
}

func (c *CoinGecko) getURL(symbol string, name string) string {
	return fmt.Sprintf("%s?%s", coinGeckoURL, c.localPairName(strings.ToLower(symbol), (strings.ToLower(name))))
}

func (c CoinGecko) Pool() query.WorkerPool {
	return c.WorkerPool
}

func (c *CoinGecko) callOne(token Token) (*CSupply, error) {
	var err error
	req := &query.HTTPRequest{
		URL: c.getURL(token.Symbol, token.Name),
		Headers: map[string]string{
			"x-cg-pro-api-key": "CG-ii7fRPp8ky22cwBm3qEianQs",
			"accept":           "application/json",
		},
	}

	res := c.Pool().Query(req)

	if res == nil {
		return nil, origins.ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return nil, res.Error
	}

	var resp []coinGeckoResponse
	err = json.Unmarshal(res.Body, &resp)

	if err != nil {
		return nil, fmt.Errorf("failed to parse CoinGecko response: %w", err)
	}

	CirculatingSupply := resp[0].CirculatingSupply
	if err != nil {
		return nil, fmt.Errorf("failed to parse CirculatingSupply float64: %w", err)
	}

	return &CSupply{
		Token:     token,
		CSupply:   CirculatingSupply,
		Timestamp: time.Now(),
	}, nil
}

func (c *CoinGecko) getCirculatingSupply(tokens []Token) []FetchResult {
	return callSinglePairOrigin(c, tokens)
}
