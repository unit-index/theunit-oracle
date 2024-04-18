package origins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

const uniswapURL = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"

type uniswapResponse struct {
	Data struct {
		Pairs []uniswapPairResponse
	}
}

type uniswapTokenResponse struct {
	Symbol string `json:"symbol"`
}

type uniswapPairResponse struct {
	ID      string               `json:"id"`
	Price0  stringAsFloat64      `json:"token0Price"`
	Price1  stringAsFloat64      `json:"token1Price"`
	Volume0 stringAsFloat64      `json:"volumeToken0"`
	Volume1 stringAsFloat64      `json:"volumeToken1"`
	Token0  uniswapTokenResponse `json:"token0"`
	Token1  uniswapTokenResponse `json:"token1"`
}

type Uniswap struct {
	WorkerPool        query.WorkerPool
	ContractAddresses ContractAddresses
}

func (u *Uniswap) pairsToContractAddresses(pairs []Pair) ([]string, error) {
	var names []string
	for _, pair := range pairs {
		address, _, ok := u.ContractAddresses.ByPair(pair)
		if !ok {
			return names, fmt.Errorf("failed to find contract address for pair %s", pair.String())
		}
		names = append(names, address)
	}
	return names, nil
}

func (u Uniswap) Pool() query.WorkerPool {
	return u.WorkerPool
}

func (u Uniswap) PullPrices(pairs []Pair) []FetchResult {
	var err error

	contracts, err := u.pairsToContractAddresses(pairs)
	if err != nil {
		return fetchResultListWithErrors(pairs, err)
	}
	pairsJSON, _ := json.Marshal(contracts)
	gql := `
		query($ids:[String]) {
			pairs(where:{id_in:$ids}) {
				id
				token0Price
				token1Price
				volumeToken0
				volumeToken1
				token0 { symbol }
				token1 { symbol }
			}
		}
	`
	body := fmt.Sprintf(
		`{"query":"%s","variables":{"ids":%s}}`,
		strings.ReplaceAll(strings.ReplaceAll(gql, "\n", " "), "\t", ""),
		pairsJSON,
	)

	req := &query.HTTPRequest{
		URL:    uniswapURL,
		Method: "POST",
		Body:   bytes.NewBuffer([]byte(body)),
	}

	// make query
	res := u.WorkerPool.Query(req)
	if res == nil {
		return fetchResultListWithErrors(pairs, ErrEmptyOriginResponse)
	}
	if res.Error != nil {
		return fetchResultListWithErrors(pairs, res.Error)
	}

	// parse JSON
	var resp uniswapResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return fetchResultListWithErrors(pairs, fmt.Errorf("failed to parse Uniswap response: %w", err))
	}

	// convert response from a slice to a map
	respMap := map[string]uniswapPairResponse{}
	for _, pairResp := range resp.Data.Pairs {
		respMap[pairResp.Token0.Symbol+"/"+pairResp.Token1.Symbol] = pairResp
	}

	// prepare result
	results := make([]FetchResult, 0)
	for _, pair := range pairs {
		b := pair.Base
		q := pair.Quote

		pair0 := b + "/" + q
		pair1 := q + "/" + b

		if r, ok := respMap[pair0]; ok {
			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     r.Price1.val(),
					Bid:       r.Price1.val(),
					Ask:       r.Price1.val(),
					Volume24h: r.Volume0.val(),
					Timestamp: time.Now(),
				},
			})
		} else if r, ok := respMap[pair1]; ok {
			results = append(results, FetchResult{
				Price: Price{
					Pair:      pair,
					Price:     r.Price0.val(),
					Bid:       r.Price0.val(),
					Ask:       r.Price0.val(),
					Volume24h: r.Volume1.val(),
					Timestamp: time.Now(),
				},
			})
		} else {
			results = append(results, FetchResult{
				Price: Price{Pair: pair},
				Error: ErrMissingResponseForPair,
			})
		}
	}

	return results
}
