package unit

import (
	"fmt"
	"strings"
	"time"
)

/*

var coinPrice = []

for price := coinPrice {
	(price / lastMonthPrice ) * lastMonthWight
}

*/

type Token struct {
	Name   string
	Symbol string
}

type CSupply struct {
	CSupply    float64
	Type       string
	Parameters map[string]string
	Time       time.Time
	Token      Token
	Error      string
	CSupplys   []*CSupply
}
type Unit interface {
	TokenTotalSupply(tokens Token) (*CSupply, error)
	TokensTotalSupply(tokens ...Token) (map[Token]*CSupply, error)
}

func calculationUNIT(tokens []Token) {
	//var count = 0
	//for _, token := range tokens {
	//	count = token.Price / token.lastMonthPrice * token.lastMonthWight
	//}
}

func NewTokens(s ...string) ([]Token, error) {
	var t []Token
	for _, p := range s {
		pr, err := NewToken(p)
		if err != nil {
			return nil, err
		}
		t = append(t, pr)
	}
	return t, nil
}

func (p Token) String() string {
	return fmt.Sprintf("%s-%s", p.Name, p.Symbol)
}

func (p Token) Equal(c Token) bool {
	return p.Name == c.Name && p.Symbol == c.Symbol
}

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
