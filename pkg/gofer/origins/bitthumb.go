package origins

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

// Exchange URL
const bitThumpURL = "https://global-openapi.bithumb.pro/openapi/v1/spot/ticker?symbol=%s"

type bitThumbPriceResponse struct {
	Low    stringAsFloat64 `json:"l"`
	High   stringAsFloat64 `json:"h"`
	Last   stringAsFloat64 `json:"c"`
	Symbol string          `json:"s"`
	Volume stringAsFloat64 `json:"v"`
}
type bitThumbResponse struct {
	Data      []bitThumbPriceResponse `json:"data"`
	Code      string                  `json:"code"`
	Msg       string                  `json:"msg"`
	Timestamp intAsUnixTimestampMs    `json:"timestamp"`
}

// Bithumb origin handler
type BitThump struct {
	WorkerPool query.WorkerPool
}

func (c *BitThump) localPairName(pair Pair) string {
	return fmt.Sprintf("%s-%s", strings.ToUpper(pair.Base), strings.ToUpper(pair.Quote))
}

func (c *BitThump) getURL(pair Pair) string {
	return fmt.Sprintf(bitThumpURL, c.localPairName(pair))
}

func (c BitThump) Pool() query.WorkerPool {
	return c.WorkerPool
}

func (c BitThump) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&c, pairs)
}

func (c *BitThump) callOne(pair Pair) (*Price, error) {
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
	var resp bitThumbResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Bithumb response: %w", err)
	}
	if resp.Code != "0" || resp.Msg != "success" || len(resp.Data) != 1 {
		return nil, fmt.Errorf("invalid Bithumb response: %s", res.Body)
	}
	priceResp := resp.Data[0]
	// building Price
	return &Price{
		Pair:      pair,
		Price:     priceResp.Last.val(),
		Volume24h: priceResp.Volume.val(),
		Timestamp: resp.Timestamp.val(),
	}, nil
}
