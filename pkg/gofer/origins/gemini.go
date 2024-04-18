package origins

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

// Gemini URL
const geminiURL = "https://api.gemini.com/v1/pubticker/%s"

type geminiResponse struct {
	Price string `json:"last"`
	Ask   string `json:"ask"`
	Bid   string `json:"bid"`
}

// Gemini origin handler
type Gemini struct {
	WorkerPool query.WorkerPool
}

func (g *Gemini) localPairName(pair Pair) string {
	return strings.ToLower(pair.Base + pair.Quote)
}

func (g *Gemini) getURL(pair Pair) string {
	return fmt.Sprintf(geminiURL, g.localPairName(pair))
}

func (g Gemini) Pool() query.WorkerPool {
	return g.WorkerPool
}

func (g Gemini) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&g, pairs)
}

func (g *Gemini) callOne(pair Pair) (*Price, error) {
	var err error
	req := &query.HTTPRequest{
		URL: g.getURL(pair),
	}

	// make query
	res := g.Pool().Query(req)
	if res == nil {
		return nil, ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return nil, res.Error
	}
	// parsing JSON
	var resp geminiResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w", err)
	}
	// Parsing price from string
	price, err := strconv.ParseFloat(resp.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price from gemini origin %s", res.Body)
	}
	// Parsing ask from string
	ask, err := strconv.ParseFloat(resp.Ask, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ask from gemini origin %s", res.Body)
	}
	// Parsing bid from string
	bid, err := strconv.ParseFloat(resp.Bid, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bid from gemini origin %s", res.Body)
	}
	// building Price
	return &Price{
		Pair:      pair,
		Price:     price,
		Ask:       ask,
		Bid:       bid,
		Timestamp: time.Now(),
	}, nil
}
