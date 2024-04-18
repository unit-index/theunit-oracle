package origins

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

// Huobi URL
const huobiURL = "https://api.huobi.pro/market/tickers"

type huobiResponse struct {
	Symbol string  `json:"symbol"`
	Volume float64 `json:"vol"`
	Bid    float64 `json:"bid"`
	Ask    float64 `json:"ask"`
}

// Huobi origin handler
type Huobi struct {
	WorkerPool query.WorkerPool
}

func (h *Huobi) localPairName(pair Pair) string {
	return strings.ToLower(pair.Base + pair.Quote)
}

func (h Huobi) Pool() query.WorkerPool {
	return h.WorkerPool
}

func (h Huobi) PullPrices(pairs []Pair) []FetchResult {
	frs, err := h.fetch(pairs)
	if err != nil {
		return fetchResultListWithErrors(pairs, err)
	}
	return frs
}

func (h *Huobi) fetch(pairs []Pair) ([]FetchResult, error) {
	var err error
	req := &query.HTTPRequest{
		URL: huobiURL,
	}

	res := h.Pool().Query(req)
	if res == nil {
		return nil, ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return nil, res.Error
	}

	var resp struct {
		Status    string          `json:"status"`
		Timestamp int64           `json:"ts"`
		Data      []huobiResponse `json:"data"`
	}

	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse huobi response: %w", err)
	}
	if resp.Status == "error" {
		return nil, fmt.Errorf("error response from huobi origin %s", res.Body)
	}

	respMap := map[string]huobiResponse{}
	for _, t := range resp.Data {
		respMap[t.Symbol] = t
	}

	ts := time.Unix(resp.Timestamp/1000, 0)
	frs := make([]FetchResult, len(pairs))
	for i, p := range pairs {
		if t, has := respMap[h.localPairName(p)]; has {
			frs[i] = fetchResult(Price{
				Pair:      p,
				Price:     (t.Ask + t.Bid) / 2,
				Ask:       t.Ask,
				Bid:       t.Bid,
				Volume24h: t.Volume,
				Timestamp: ts,
			})
		} else {
			frs[i] = fetchResultWithError(
				p,
				fmt.Errorf("failed to find symbol %s in huobi response", p),
			)
		}
	}

	return frs, nil
}
