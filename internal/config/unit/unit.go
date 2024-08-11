package unit

import (
	"context"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/internal/query"
	pkgEthereum "github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/feeder"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/nodes"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/origins"
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

type Token struct {
	Name                     string   `json:"name"`
	Symbol                   string   `json:"symbol"`
	Method                   string   `json:"method"`
	MinimumSuccessfulSources int      `json:"minimumSuccessfulSources"`
	CirculatingSupplySource  []string `json:"circulatingSupplySource"`
}

type Unit struct {
	Tokens                  []Token                   `json:"tokens"`
	CirculatingSupplySource []CirculatingSupplySource `json:"circulatingSupplySource"`
}

func (u *Unit) Configure() {

}

func (u *Unit) TokenTotalSupply(tokens []unit.Token) {

}

func (u *Unit) ConfigureUnit(ctx context.Context, cli pkgEthereum.Client, logger log.Logger, noRPC bool) (unit.Unit, error) {
	gra, err := u.buildGraphs()
	originSet, err := u.buildOrigins()
	if err != nil {
		return nil, err
	}
	fed := feeder.NewFeeder(ctx, originSet, logger)

	unit := graph.NewUnit(gra, fed)
	return unit, nil
}

func (u *Unit) buildOrigins() (*origins.Set, error) {
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

func (u *Unit) buildGraphs() (map[unit.Token]nodes.Aggregator, error) {
	var err error

	graphs := map[unit.Token]nodes.Aggregator{}

	// It's important to create root nodes before branches, because branches
	// may refer to another root nodes instances.
	err = u.buildRoots(graphs)
	if err != nil {
		return nil, err
	}

	//err = u.buildBranches(graphs)
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = u.detectCycle(graphs)
	//if err != nil {
	//	return nil, err
	//}

	return graphs, nil
}

func (u *Unit) buildRoots(graphs map[unit.Token]nodes.Aggregator) error {
	for _, model := range u.Tokens {
		modelToken, err := unit.NewToken(model.Name + ":" + model.Symbol)
		if err != nil {
			return err
		}

		switch model.Method {
		case "median":
			graphs[modelToken] = nodes.NewMedianAggregatorNode(modelToken, model.MinimumSuccessfulSources)
		default:
			return fmt.Errorf("unknown method %s for pair %s", model.Method, model.Name)
		}
	}

	return nil
}
