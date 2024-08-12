package origins

import (
	"errors"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/internal/query"
	"sync"
	"time"
)

var ErrUnknownOrigin = errors.New("unknown origin")

type Handler interface {
	// Fetch should implement making API request to origin URL and
	// collecting/parsing origin data.
	Fetch(tokens []Token) []FetchResult
}

type ExchangeHandler interface {
	getCirculatingSupply(tokens []Token) []FetchResult
}

type CSupply struct {
	Token     Token
	CSupply   float64
	Timestamp time.Time
}

type FetchResult struct {
	CSupply CSupply
	Error   error
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

func (p Token) Equal(c Token) bool {
	return p.Name == c.Name && p.Symbol == c.Symbol
}

type Token struct {
	Name   string
	Symbol string
}

func NewBaseHandler(handler ExchangeHandler) ExchangeHandler {
	return &BaseExchangeHandler{ExchangeHandler: handler}
}

func NewBaseExchangeHandler(handler ExchangeHandler) *BaseExchangeHandler {
	return &BaseExchangeHandler{
		ExchangeHandler: handler,
	}
}

func (h BaseExchangeHandler) Fetch(tokens []Token) []FetchResult {
	//if h.aliases == nil {
	//	return h.PullPrices(pairs)
	//}
	//
	//var renamedPairs []Pair
	//for _, token := range tokens {
	//	renamedPairs = append(renamedPairs, h.aliases.replacePair(pair))
	//}
	results := h.getCirculatingSupply(tokens)

	//// Reverting our replacement
	//for i := range results {
	//	results[i].Price.Pair = h.aliases.revertPair(results[i].Price.Pair)
	//}
	return results
	//return []FetchResult{}
}

func (e *Set) Fetch(originTokens map[string][]Token) map[string][]FetchResult {
	var mu sync.Mutex
	var wg sync.WaitGroup
	ch := make(chan struct{}, e.goroutines)

	wg.Add(len(originTokens))
	frs := map[string][]FetchResult{}
	for origin, tokens := range originTokens {
		ch <- struct{}{}

		origin, tokens := origin, tokens
		handler, ok := e.list[origin]

		go func() {
			defer func() { <-ch }()

			if !ok {
				mu.Lock()
				frs[origin] = fetchResultListWithErrors(
					tokens,
					fmt.Errorf("%w (%s)", ErrUnknownOrigin, origin),
				)
				mu.Unlock()
			} else {
				resp := handler.Fetch(tokens)
				mu.Lock()
				frs[origin] = append(frs[origin], resp...)
				mu.Unlock()
			}

			wg.Done()
		}()
	}

	wg.Wait()
	return frs
}

func fetchResultListWithErrors(tokens []Token, err error) []FetchResult {
	r := make([]FetchResult, len(tokens))
	for i, token := range tokens {
		r[i] = FetchResult{
			CSupply: CSupply{
				Token:     token,
				Timestamp: time.Now(),
			},
			Error: err,
		}
	}
	return r
}

func NewSet(list map[string]Handler, goroutines int) *Set {
	return &Set{list: list, goroutines: goroutines}
}

func DefaultOriginSet(_ query.WorkerPool, goroutines int) *Set {
	return NewSet(map[string]Handler{}, goroutines)
}

type singleTokenOrigin interface {
	callOne(token Token) (*CSupply, error)
}

func callSinglePairOrigin(e singleTokenOrigin, tokens []Token) []FetchResult {
	crs := make([]FetchResult, 0)
	for _, token := range tokens {
		cSupply, err := e.callOne(token)
		if err != nil {
			crs = append(crs, FetchResult{
				CSupply: CSupply{Token: token},
				Error:   err,
			})
		} else {
			crs = append(crs, FetchResult{
				CSupply: *cSupply,
				Error:   err,
			})
		}
	}

	return crs
}
