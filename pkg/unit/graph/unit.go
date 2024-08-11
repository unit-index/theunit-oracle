package graph

import (
	"fmt"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/feeder"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/nodes"
)

type ErrTokenNotFound struct {
	Token unit.Token
}

func (e ErrTokenNotFound) Error() string {
	return fmt.Sprintf("unable to find the %s pair", e.Token)
}

type Unit struct {
	graphs map[unit.Token]nodes.Aggregator
	feeder *feeder.Feeder
}

func NewUnit(g map[unit.Token]nodes.Aggregator, f *feeder.Feeder) *Unit {
	return &Unit{graphs: g, feeder: f}
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
