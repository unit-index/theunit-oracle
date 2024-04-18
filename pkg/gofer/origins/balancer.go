package origins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/toknowwhy/theunit-oracle/internal/query"
)

const balancerURL = "https://api.thegraph.com/subgraphs/name/balancer-labs/balancer"

type balancerResponse struct {
	Data struct {
		TokenPrices []balancerPairResponse `json:"tokenPrices"`
	}
}

type balancerPairResponse struct {
	Symbol string          `json:"symbol"`
	Price  stringAsFloat64 `json:"price"`
	Volume stringAsFloat64 `json:"poolLiquidity"`
}

type Balancer struct {
	WorkerPool        query.WorkerPool
	ContractAddresses ContractAddresses
}

func (s Balancer) Pool() query.WorkerPool {
	return s.WorkerPool
}

func (s *Balancer) pairsToContractAddress(pair Pair) (string, error) {
	// We're checking for reverse pairs because the same contract is used to
	// trade in both directions.
	contract, _, ok := s.ContractAddresses.ByPair(pair)
	if !ok {
		return "", fmt.Errorf("failed to find Balancer contract address for pair: %s", pair.String())
	}
	return contract, nil
}

func (s Balancer) PullPrices(pairs []Pair) []FetchResult {
	return callSinglePairOrigin(&s, pairs)
}

func (s *Balancer) callOne(pair Pair) (*Price, error) {
	var err error

	contract, err := s.pairsToContractAddress(pair)
	if err != nil {
		return nil, err
	}
	pairsJSON, _ := json.Marshal(contract)
	gql := `
		query($id:String) {
			tokenPrices(where: { id: $id }){
				symbol price poolLiquidity
			}
		}
	`
	body := fmt.Sprintf(
		`{"query":"%s","variables":{"id":%s}}`,
		strings.ReplaceAll(strings.ReplaceAll(gql, "\n", " "), "\t", ""),
		pairsJSON,
	)

	req := &query.HTTPRequest{
		URL:    balancerURL,
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
	var resp balancerResponse
	err = json.Unmarshal(res.Body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Balancer response: %w", err)
	}

	pairPrice := resp.Data.TokenPrices[0]
	if pairPrice.Symbol != pair.Base {
		return nil, ErrMissingResponseForPair
	}

	return &Price{
		Pair:      pair,
		Price:     pairPrice.Price.val(),
		Volume24h: pairPrice.Volume.val(),
		Timestamp: time.Now(),
	}, nil
}
