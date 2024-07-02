package nodes

import (
	"fmt"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

type OriginToken struct {
	Origin string
	Pair   gofer.Pair
}

func (o OriginToken) String() string {
	return fmt.Sprintf("%s %s", o.Pair.String(), o.Origin)
}

// OriginPrice represent a price which was sourced directly from an origin.
type OriginSupply struct {
	Supply string
	Origin string
	Error  error
}
