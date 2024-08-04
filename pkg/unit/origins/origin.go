package origins

import (
	"github.com/toknowwhy/theunit-oracle/internal/query"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
)

type Handler interface {
	// Fetch should implement making API request to origin URL and
	// collecting/parsing origin data.
	Fetch(tokens []unit.Token) []float64
}

type ExchangeHandler interface {
	getCirculatingSupply(symbol string, name string) (string, error)
}

type BaseExchangeHandler struct {
	ExchangeHandler
}

type Set struct {
	list       map[string]Handler
	goroutines int
}

func (e *Set) SetHandler(name string, handler Handler) {
	e.list[name] = handler
}

func NewBaseHandler(handler ExchangeHandler) Handler {
	return &BaseExchangeHandler{ExchangeHandler: handler}
}

func NewBaseExchangeHandler(handler ExchangeHandler) *BaseExchangeHandler {
	return &BaseExchangeHandler{
		ExchangeHandler: handler,
	}
}

func (h BaseExchangeHandler) Fetch(tokens []unit.Token) []float64 {
	//if h.aliases == nil {
	//	return h.PullPrices(pairs)
	//}
	//
	//var renamedPairs []Pair
	//for _, pair := range pairs {
	//	renamedPairs = append(renamedPairs, h.aliases.replacePair(pair))
	//}
	//results := h.PullPrices(renamedPairs)
	//
	//// Reverting our replacement
	//for i := range results {
	//	results[i].Price.Pair = h.aliases.revertPair(results[i].Price.Pair)
	//}
	//return results
	return []float64{}
}

func NewSet(list map[string]Handler, goroutines int) *Set {
	return &Set{list: list, goroutines: goroutines}
}

func DefaultOriginSet(_ query.WorkerPool, goroutines int) *Set {
	return NewSet(map[string]Handler{}, goroutines)
}
