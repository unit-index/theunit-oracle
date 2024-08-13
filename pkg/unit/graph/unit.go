package graph

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/oracle"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/feeder"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/nodes"
	"math/big"
	"time"
)

type ErrTokenNotFound struct {
	Token unit.Token
}

func (e ErrTokenNotFound) Error() string {
	return fmt.Sprintf("unable to find the %s pair", e.Token)
}

type Unit struct {
	graphs        map[unit.Token]nodes.Aggregator
	feeder        *feeder.Feeder
	unitAlgorithm oracle.UnitAlgorithm
	gofer         gofer.Gofer
	tokens        map[common.Address]unit.Token
}

func NewUnit(g map[unit.Token]nodes.Aggregator, f *feeder.Feeder, unitAlgorithm oracle.UnitAlgorithm, gofer gofer.Gofer, tokens map[common.Address]unit.Token) *Unit {
	return &Unit{graphs: g, feeder: f, unitAlgorithm: unitAlgorithm, gofer: gofer, tokens: tokens}
}

func (u *Unit) TokenTotalSupply(token unit.Token) (*unit.CSupply, error) {
	n, ok := u.graphs[token]
	if !ok {
		return nil, ErrTokenNotFound{Token: token}
	}

	if u.feeder != nil {
		u.feeder.Feed(n)
	}
	return mapGraphCSupply(n.CSupply()), nil
}

func (u *Unit) TokensTotalSupply(tokens ...unit.Token) (map[unit.Token]*unit.CSupply, error) {
	ns, err := u.findNodes(tokens...)
	if err != nil {
		return nil, err
	}
	if u.feeder != nil {
		u.feeder.Feed(ns...)
	}
	res := make(map[unit.Token]*unit.CSupply)
	for _, n := range ns {
		if n, ok := n.(nodes.Aggregator); ok {
			res[n.Token()] = mapGraphCSupply(n.CSupply())
		}
	}
	return res, nil
}

func (u *Unit) FeedMarketCapAndPrice(tokens ...unit.Token) ([]unit.UnitPerMonthParams, error) {
	css, err := u.TokensTotalSupply(tokens...)

	var unitPerMonthParams []unit.UnitPerMonthParams

	//  去合约里查本月有哪些token在index内
	//获取上个月最后一天的时间戳。然后再去合约里查
	now := time.Now()
	firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfLastMonth := firstDayOfThisMonth.Add(-time.Second)
	lastDayOfLastMonthUTC := lastDayOfLastMonth.In(time.UTC)

	// 获取时间戳
	timestamp := lastDayOfLastMonthUTC.Unix()
	fmt.Println("FeedMarketCapAndPrice", timestamp)
	bigIntTimestamp := big.NewInt(0).SetInt64(timestamp)
	ctx := context.Background()
	token, err := u.unitAlgorithm.GetTokens(ctx, bigIntTimestamp)
	if err != nil {
		return nil, err
	}
	fmt.Println(token)

	// 去调合约
	for _, token := range u.tokens {

		//csupply, err := u.TokenTotalSupply(token)
		//fmt.Println("marketCap", token.Symbol)
		if err != nil {
			return nil, err
		}
		pair, err := gofer.NewPair(fmt.Sprintf("%s/%s", token.Symbol, "USD"))
		if err != nil {
			return nil, err
		}
		p, err := u.gofer.Price(pair)
		if err != nil {
			return nil, err
		}

		CSupply := css[token].CSupply

		marketCap := CSupply * p.Price

		fmt.Println("marketCap", marketCap)

		unitPerMonthParams = append(unitPerMonthParams, unit.UnitPerMonthParams{CSupply: CSupply, LastPrice: p.Price, LastMarketCap: marketCap})
	}

	return unitPerMonthParams, nil
}

func (u *Unit) Price() (string, error) {
	return "", nil
}

func mapGraphCSupply(t interface{}) *unit.CSupply {
	gt := &unit.CSupply{
		Parameters: make(map[string]string),
	}

	switch typedCSupply := t.(type) {
	case nodes.AggregatorCSupply:
		gt.Type = "aggregator"
		gt.Token = typedCSupply.Token
		gt.CSupply = typedCSupply.CSupply
		gt.Time = typedCSupply.Time
		if typedCSupply.Error != nil {
			gt.Error = typedCSupply.Error.Error()
		}
		gt.Parameters = typedCSupply.Parameters
		for _, ct := range typedCSupply.OriginCSupply {
			gt.CSupplys = append(gt.CSupplys, mapGraphCSupply(ct))
		}
		for _, ct := range typedCSupply.AggregatorPrices {
			gt.CSupplys = append(gt.CSupplys, mapGraphCSupply(ct))
		}
	case nodes.OriginCSupply:
		gt.Type = "origin"
		gt.Token = typedCSupply.Token
		gt.CSupply = typedCSupply.CSupply
		gt.Time = typedCSupply.Time
		if typedCSupply.Error != nil {
			gt.Error = typedCSupply.Error.Error()
		}
		gt.Parameters["origin"] = typedCSupply.Origin
	default:
		panic("unsupported object")
	}

	return gt
}

func (u *Unit) findNodes(tokens ...unit.Token) ([]nodes.Node, error) {
	var ns []nodes.Node
	if len(tokens) == 0 { // Return all:
		for _, n := range u.graphs {
			ns = append(ns, n)
		}
	} else { // Return for given pairs:
		for _, p := range tokens {
			n, ok := u.graphs[p]
			if !ok {
				return nil, ErrTokenNotFound{Token: p}
			}
			ns = append(ns, n)
		}
	}
	return ns, nil
}
