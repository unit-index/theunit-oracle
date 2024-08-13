package unit

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/toknowwhy/theunit-oracle/internal/query"
	pkgEthereum "github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/oracle/geth"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/feeder"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/nodes"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/origins"
	"reflect"
	"sort"
	"strings"
	"time"
)

//type Token struct {
//	Name           string `json:"name"`
//	Symbol         string `json:"symbol"`
//	Price          float64
//	lastMonthPrice float64
//	lastMonthWight float64
//}

const defaultTTL = 60 * time.Second
const maxTTL = 60 * time.Second

type ErrCyclicReference struct {
	Token unit.Token
	Path  []nodes.Node
}

func (e ErrCyclicReference) Error() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("a cyclic reference was detected for the %s pair: ", e.Path))
	for i, n := range e.Path {
		t := reflect.TypeOf(n).String()
		switch typedNode := n.(type) {
		case nodes.Aggregator:
			s.WriteString(fmt.Sprintf("%s(%s)", t, typedNode.Token()))
		default:
			s.WriteString(t)
		}
		if i != len(e.Path)-1 {
			s.WriteString(" -> ")
		}
	}
	return s.String()
}

type CirculatingSupplySource struct {
	Origin string `json:"origin"`
	Key    string `json:"key"`
}

type Token struct {
	Name                     string   `json:"name"`
	Symbol                   string   `json:"symbol"`
	Method                   string   `json:"method"`
	Address                  string   `json:"address"`
	MinimumSuccessfulSources int      `json:"minimumSuccessfulSources"`
	CirculatingSupplySource  []string `json:"circulatingSupplySource"`
	TTL                      int      `json:"ttl"`
}

type Unit struct {
	Tokens                  []Token                   `json:"tokens"`
	CirculatingSupplySource []CirculatingSupplySource `json:"circulatingSupplySource"`
	FeedAddress             string                    `json:"feedAddress"`
}

func (u *Unit) Configure() {

}

func (u *Unit) TokenTotalSupply(tokens []unit.Token) {

}

func (u *Unit) ConfigureUnit(ctx context.Context, cli pkgEthereum.Client, gofer gofer.Gofer, logger log.Logger, noRPC bool) (unit.Unit, error) {
	gra, err := u.buildGraphs()
	originSet, err := u.buildOrigins()
	if err != nil {
		return nil, err
	}
	fed := feeder.NewFeeder(ctx, originSet, logger)

	feedAddress := pkgEthereum.HexToAddress(u.FeedAddress)
	unitAlgorithm := geth.NewUnitAlgorithm(cli, feedAddress)

	var tokens = make(map[common.Address]unit.Token)
	for _, token := range u.Tokens {
		//fmt.Println(token.Address, common.HexToAddress(token.Address))
		tokens[common.HexToAddress(token.Address)] = unit.Token{Name: token.Name, Symbol: token.Symbol}
	}

	unit := graph.NewUnit(gra, fed, unitAlgorithm, gofer, tokens)
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

	err = u.buildBranches(graphs)
	if err != nil {
		return nil, err
	}

	err = u.detectCycle(graphs)
	if err != nil {
		return nil, err
	}

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

func (c *Unit) buildBranches(graphs map[unit.Token]nodes.Aggregator) error {
	for _, model := range c.Tokens {
		// We can ignore error here, because it was checked already
		// in buildRoots method.
		modelToken, err := unit.NewToken(model.Name + ":" + model.Symbol)
		if err != nil {
			return err
		}

		var parent nodes.Parent
		if typedNode, ok := graphs[modelToken].(nodes.Parent); ok {
			parent = typedNode
		} else {
			return fmt.Errorf(
				"%s must implement the nodes.Parent interface",
				reflect.TypeOf(graphs[modelToken]).Elem().String(),
			)
		}

		for _, source := range model.CirculatingSupplySource {
			//var children []nodes.Node
			//for _, source := range sources {
			//	var err error
			var node nodes.Node
			//
			//	if source.Origin == "." {
			//		node, err = c.reference(graphs, source)
			//		if err != nil {
			//			return err
			//		}
			//	} else {
			//		node, err = c.originNode(model, source)
			//		if err != nil {
			//			return err
			//		}
			//	}
			//
			//	children = append(children, node)
			//}
			node, err = c.originNode(model, source)
			if err != nil {
				return err
			}
			// If there are provided multiple sources it means, that the price
			// have to be calculated by using the nodes.IndirectAggregatorNode.
			// Otherwise we can pass that nodes.OriginNode directly to
			// the parent node.
			//var node nodes.Node
			//if len(children) == 1 {
			//	node = children[0]
			//} else {
			//	indirectAggregator := nodes.NewIndirectAggregatorNode(modelPair)
			//	for _, c := range children {
			//		indirectAggregator.AddChild(c)
			//	}
			//	node = indirectAggregator
			//}
			//
			parent.AddChild(node)
		}
	}

	return nil
}

func (c *Unit) originNode(model Token, source string) (nodes.Node, error) {
	sourceToken, err := unit.NewToken(model.Name + ":" + model.Symbol)
	if err != nil {
		return nil, err
	}

	originPair := nodes.OriginToken{
		Origin: source,
		Token:  sourceToken,
	}

	ttl := defaultTTL
	if model.TTL > 0 {
		ttl = time.Second * time.Duration(model.TTL)
	}
	//if source.TTL > 0 {
	//	ttl = time.Second * time.Duration(source.TTL)
	//}

	return nodes.NewOriginNode(originPair, ttl, ttl+maxTTL), nil
}
func (c *Unit) detectCycle(graphs map[unit.Token]nodes.Aggregator) error {
	for _, token := range sortGraphs(graphs) {
		if path := nodes.DetectCycle(graphs[token]); len(path) > 0 {
			return ErrCyclicReference{Token: token, Path: path}
		}
	}

	return nil
}

func sortGraphs(graphs map[unit.Token]nodes.Aggregator) []unit.Token {
	var ps []unit.Token
	for p := range graphs {
		ps = append(ps, p)
	}
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].String() < ps[j].String()
	})
	return ps
}
