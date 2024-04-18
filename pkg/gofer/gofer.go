package gofer

import (
	"fmt"
	"strings"
	"time"
)

// Pair represents an asset pair.
type Pair struct {
	Base  string
	Quote string
}

type Token struct {
	Symbol string
}

// NewPair returns a new Pair for given string. The string must be formatted
// as "BASE/QUOTE".
func NewPair(s string) (Pair, error) {
	ss := strings.Split(s, "/")
	if len(ss) != 2 {
		return Pair{}, fmt.Errorf("couldn't parse pair \"%s\"", s)
	}
	return Pair{Base: strings.ToUpper(ss[0]), Quote: strings.ToUpper(ss[1])}, nil
}

// NewPairs returns a Pair slice for given strings. Given strings must be
// formatted as "BASE/QUOTE".
func NewPairs(s ...string) ([]Pair, error) {
	var r []Pair
	for _, p := range s {
		pr, err := NewPair(p)
		if err != nil {
			return nil, err
		}
		r = append(r, pr)
	}
	return r, nil
}

func NewToken(s ...string) ([]Token, error) {
	var t []Token
	for _, p := range s {
		to := Token{p}
		t = append(t, to)
	}
	return t, nil
}

func (p Pair) Empty() bool {
	return p.Base == "" && p.Quote == ""
}

func (p Pair) Equal(c Pair) bool {
	return p.Base == c.Base && p.Quote == c.Quote
}

func (p Pair) String() string {
	return fmt.Sprintf("%s/%s", p.Base, p.Quote)
}

// Model is a simplified representation of a model which is used to calculate
// asset pair prices. The main purpose of this structure is to help the end
// user to understand how prices are derived and calculated.
//
// This structure is purely informational. The way it is used depends on
// a specific implementation.
type Model struct {
	// Type is used to differentiate between model types.
	Type string
	// Parameters is a optional list of model's parameters.
	Parameters map[string]string
	// Pair is a asset pair for which this model returns a price.
	Pair Pair
	// Models is a list of sub models used to calculate price.
	Models []*Model
}

// Price represents price for a single pair. If the Price price was calculated
// indirectly it will also contain all prices used to calculate the price.
type Price struct {
	Type       string
	Parameters map[string]string
	Pair       Pair
	Price      float64
	Bid        float64
	Ask        float64
	Volume24h  float64
	Time       time.Time
	Prices     []*Price
	Error      string
}

// Gofer provides prices for asset pairs.
type Gofer interface {
	// Models describes price models which are used to calculate prices.
	// If no pairs are specified, models for all pairs are returned.
	Models(pairs ...Pair) (map[Pair]*Model, error)
	// Price returns a Price for the given pair.
	Price(pair Pair) (*Price, error)
	// Prices returns prices for the given pairs. If no pairs are specified,
	// prices for all pairs are returned.
	Prices(pairs ...Pair) (map[Pair]*Price, error)
	// Pairs returns all pairs.
	Pairs() ([]Pair, error)

	//Supply(token Token) (float64, error)
}

// StartableGofer interface represents a Gofer instances that have to be
// started first to work properly.
type StartableGofer interface {
	Gofer
	Start() error
	Wait()
}
