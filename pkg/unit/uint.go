package unit

import (
	"fmt"
	"strings"
)

/*

var coinPrice = []

for price := coinPrice {
	(price / lastMonthPrice ) * lastMonthWight
}

*/

type Token struct {
	Name           string
	Symbol         string
	Price          float64
	lastMonthPrice float64
	lastMonthWight float64
}
type Unit interface {
	//TokenTotalSupply([]Token)
}

func calculationUNIT(tokens []Token) {
	//var count = 0
	//for _, token := range tokens {
	//	count = token.Price / token.lastMonthPrice * token.lastMonthWight
	//}
}

func NewTokens(s ...string) ([]Token, error) {
	var t []Token
	for _, symbol := range s {
		to := Token{Symbol: symbol}
		t = append(t, to)
	}
	return t, nil
}

//func NewPair(s string) (Pair, error) {
//	ss := strings.Split(s, "/")
//	if len(ss) != 2 {
//		return Pair{}, fmt.Errorf("couldn't parse pair \"%s\"", s)
//	}
//	return Pair{Base: strings.ToUpper(ss[0]), Quote: strings.ToUpper(ss[1])}, nil
//}

func NewToken(s string) (Token, error) {
	ss := strings.Split(s, ":")
	if len(ss) != 2 {
		return Token{}, fmt.Errorf("couldn't parse token name and symbol \"%s\"", s)
	}
	return Token{Name: ss[0], Symbol: ss[1]}, nil
}

type StartableUnit interface {
	Unit
	Start() error
	Wait()
}
