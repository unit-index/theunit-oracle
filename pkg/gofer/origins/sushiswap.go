package origins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

const sushiswapURL = "https://api.thegraph.com/subgraphs/name/zippoxer/sushiswap-subgraph-fork"

type sushiswapResponse struct {
	Data struct {
		Pairs []sushiswapPairResponse
	}
}

type sushiswapTokenResponse struct {
	Symbol string `json:"symbol"`
}

type sushiswapPairResponse struct {
	ID      string                 `json:"id"`
	Price0  stringAsFloat64        `json:"token0Price"`
	Price1  stringAsFloat64        `json:"token1Price"`
	Volume0 stringAsFloat64        `json:"volumeToken0"`
	Volume1 stringAsFloat64        `json:"volumeToken1"`
	Token0  sushiswapTokenResponse `json:"token0"`
	Token1  sushiswapTokenResponse `json:"token1"`
}

type Sushiswap struct {
	WorkerPool        query.WorkerPool
	ContractAddresses ContractAddresses
}

func (s *Sushiswap) pairsToContractAddress(pair Pair) (string, error) {
	contract, _, ok := s.ContractAddresses.ByPair(pair)
	if !ok {
		return "", fmt.Errorf("failed to get SushiSwap contract address for pair: %s", pair.String())
	}
	return contract, nil
}

func (s Sushiswap) Pool() query.WorkerPool {
	return s.WorkerPool
}

func (s Sushiswap) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&s, pairs)
}

func (s *Sushiswap) callOne(pair Pair) (*Price, error) {
	var err error

	contract, err := s.pairsToContractAddress(pair)
	if err != nil {
		return nil, err
	}
	pairsJSON, _ := json.Marshal(contract)
	gql := `
		query($id:String) {
			pairs(where:{id: $id}) {
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
		`{"query":"%s","variables":{"id":%s}}`,
		strings.ReplaceAll(strings.ReplaceAll(gql, "\n", " "), "\t", ""),
		pairsJSON,
	)

	req := &query.HTTPRequest{
		URL:    sushiswapURL,
		Method: "POST",
		Body:   bytes.NewBuffer([]byte(body)),
	}

	// make query
	res := s.Pool().Query(req)
	if res == nil {
		return nil, ErrEmptyOriginResponse
	}
	if res.Error != nil {
		return nil, res.Error
	}

	// parse JSON
	var resp sushiswapResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Sushiswap response: %w", err)
	}

	// convert response from a slice to a map
	respMap := map[string]sushiswapPairResponse{}
	for _, pairResp := range resp.Data.Pairs {
		respMap[pairResp.Token0.Symbol+"/"+pairResp.Token1.Symbol] = pairResp
	}

	b := pair.Base
	q := pair.Quote

	pair0 := b + "/" + q
	pair1 := q + "/" + b

	if r, ok := respMap[pair0]; ok {
		return &Price{
			Pair:      pair,
			Price:     r.Price1.val(),
			Bid:       r.Price1.val(),
			Ask:       r.Price1.val(),
			Volume24h: r.Volume0.val(),
			Timestamp: time.Now(),
		}, nil
	} else if r, ok := respMap[pair1]; ok {
		return &Price{
			Pair:      pair,
			Price:     r.Price0.val(),
			Bid:       r.Price0.val(),
			Ask:       r.Price0.val(),
			Volume24h: r.Volume1.val(),
			Timestamp: time.Now(),
		}, nil
	}

	return nil, ErrMissingResponseForPair
}
