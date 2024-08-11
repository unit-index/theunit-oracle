//  Copyright (C) 2020 Maker Ecosystem Growth Holdings, INC.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as
//  published by the Free Software Foundation, either version 3 of the
//  License, or (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package nodes

import (
	"fmt"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
)

type ErrNotEnoughSources struct {
	Given int
	Min   int
}

func (e ErrNotEnoughSources) Error() string {
	return fmt.Sprintf(
		"not enough sources to calculate median, %d given but at least %d required",
		e.Given,
		e.Min,
	)
}

type ErrIncompatiblePairs struct {
	Given    unit.Token
	Expected unit.Token
}

func (e ErrIncompatiblePairs) Error() string {
	return fmt.Sprintf(
		"unable to calculate median for different pairs, %s given but %s was expected",
		e.Given,
		e.Expected,
	)
}

// MedianAggregatorNode gets Prices from all of its children and calculates
// median price.
//
//	                         -- [Origin A/B]
//	                        /
//	[MedianAggregatorNode] ---- [Origin A/B]       -- ...
//	                        \                     /
//	                         -- [AggregatorNode A/B] ---- ...
//	                                              \
//	                                               -- ...
//
// All children of this node must return a Price for the same pair.
type MedianAggregatorNode struct {
	token      unit.Token
	minSources int
	children   []Node
}

func NewMedianAggregatorNode(token unit.Token, minSources int) *MedianAggregatorNode {
	return &MedianAggregatorNode{
		token:      token,
		minSources: minSources,
	}
}

// Children implements the Node interface.
func (n *MedianAggregatorNode) Children() []Node {
	return n.children
}

// AddChild implements the Parent interface.
func (n *MedianAggregatorNode) AddChild(node Node) {
	n.children = append(n.children, node)
}

func (n *MedianAggregatorNode) Token() unit.Token {
	return n.token
}

func (n *MedianAggregatorNode) CSupply() AggregatorCSupply {
	var ts time.Time
	var csupplys []float64
	var originCSupplys []OriginCSupply
	var aggregatorPrices []AggregatorCSupply
	var err error

	for i, c := range n.children {
		// There is no need to copy errors from prices to the MedianAggregatorNode
		// because there may be enough remaining prices to calculate median price.

		var tokenCSupply TokenCSupply
		switch typedNode := c.(type) {
		case Origin:
			originCSUpply := typedNode.CSupply()
			originCSupplys = append(originCSupplys, originCSUpply)
			tokenCSupply = originCSUpply.TokenCSupply
			if originCSUpply.Error != nil {
				continue
			}
		case Aggregator:
			aggregatorPrice := typedNode.CSupply()
			aggregatorPrices = append(aggregatorPrices, aggregatorPrice)
			tokenCSupply = aggregatorPrice.TokenCSupply
			if aggregatorPrice.Error != nil {
				continue
			}
		}

		if !n.token.Equal(tokenCSupply.Token) {
			err = multierror.Append(
				err,
				ErrIncompatiblePairs{Given: tokenCSupply.Token, Expected: n.token},
			)
			continue
		}

		if tokenCSupply.CSupply > 0 {
			csupplys = append(csupplys, tokenCSupply.CSupply)
		}

		if i == 0 || tokenCSupply.Time.Before(ts) {
			ts = tokenCSupply.Time
		}
	}

	if len(csupplys) < n.minSources {
		err = multierror.Append(
			err,
			ErrNotEnoughSources{Given: len(csupplys), Min: n.minSources},
		)
	}

	return AggregatorCSupply{
		TokenCSupply: TokenCSupply{
			Token:   n.token,
			CSupply: median(csupplys),
			Time:    ts,
		},
		OriginCSupply:    originCSupplys,
		AggregatorPrices: aggregatorPrices,
		Parameters:       map[string]string{"method": "median", "minimumSuccessfulSources": strconv.Itoa(n.minSources)},
		Error:            err,
	}
}

func median(xs []float64) float64 {
	count := len(xs)
	if count == 0 {
		return 0
	}

	sort.Float64s(xs)
	if count%2 == 0 {
		m := count / 2
		x1 := xs[m-1]
		x2 := xs[m]
		return (x1 + x2) / 2
	}

	return xs[(count-1)/2]
}
