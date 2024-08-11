package nodes

import (
	"fmt"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"time"
)

type OriginToken struct {
	Origin string
	Token  unit.Token
}

func (o OriginToken) String() string {
	return fmt.Sprintf("%s %s", o.Token.String(), o.Origin)
}

type TokenCSupply struct {
	Token   unit.Token
	CSupply float64
	Time    time.Time
}

// OriginPrice represent a price which was sourced directly from an origin.
type OriginCSupply struct {
	TokenCSupply
	// Origin is a name of Price source.
	Origin string
	// Error is a list of optional error messages which may occur during
	// calculating the price. If this string is not empty, then the price
	// value is not reliable.
	Error error
}

// AggregatorPrice represent a price which was calculated by using other prices.
type AggregatorCSupply struct {
	TokenCSupply
	// OriginPrices is a list of all OriginPrices used to calculate Price.
	OriginCSupply []OriginCSupply
	// AggregatorPrices is a list of all OriginPrices used to calculate Price.
	AggregatorPrices []AggregatorCSupply
	// Parameters is a custom list of optional parameters returned by an aggregator.
	Parameters map[string]string
	// Errors is a list of optional error messages which may occur during
	// fetching Price. If this list is not empty, then the price value
	// is not reliable.
	Error error
}
