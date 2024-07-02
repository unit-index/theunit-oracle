package unit

import (
	"context"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/internal/query"
	pkgEthereum "github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/graph/feeder"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/origins"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	pkgUnit "github.com/toknowwhy/theunit-oracle/pkg/unit"
)

//type Token struct {
//	Name           string `json:"name"`
//	Symbol         string `json:"symbol"`
//	Price          float64
//	lastMonthPrice float64
//	lastMonthWight float64
//}

type CirculatingSupplySource struct {
	Origin string `json:"origin"`
	Key    string `json:"key"`
}

type Unit struct {
	Tokens                  []pkgUnit.Token
	CirculatingSupplySource []CirculatingSupplySource `json:"circulatingSupplySource"`
}

func (u *Unit) Configure() {

}

func (u *Unit) TokenTotalSupply(tokens []pkgUnit.Token) {

}

func (u *Unit) ConfigureUnit(ctx context.Context, cli pkgEthereum.Client, logger log.Logger, noRPC bool) (Unit, error) {

	originSet, err := u.buildOrigins(cli)

	fed := feeder.NewFeeder(ctx, originSet, logger)

	return Unit{}, nil
}

func (u *Unit) buildOrigins(cli pkgEthereum.Client) (*origins.Set, error) {
	const defaultWorkerCount = 5
	wp := query.NewHTTPWorkerPool(defaultWorkerCount)
	originSet := origins.DefaultOriginSet(wp, defaultWorkerCount)
	for _, origin := range u.CirculatingSupplySource {
		handler, err := NewHandler(origin.Origin, wp, origin.Key)
		if err != nil || handler == nil {
			return nil, fmt.Errorf("failed to initiate %s origin with name %s due to error: %w",
				origin.Origin, origin.Key, err)
		}
		originSet.SetHandler(origin.Origin, handler)
	}
	return originSet, nil
}
