package origins

import "github.com/toknowwhy/theunit-oracle/pkg/unit"

type Handler interface {
	// Fetch should implement making API request to origin URL and
	// collecting/parsing origin data.
	Fetch(tokens []unit.Token) []float64
}

type ExchangeHandler interface {
	//getCirculatingSupply(symbol string, name string) (string, error)
}

func NewBaseHandler(handler ExchangeHandler) Handler {

}
