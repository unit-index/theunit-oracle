package supply

import "context"

type FetchResult struct {
	Circulating string
}
type Handler interface {
	Fetch(token string) FetchResult
}

type TotalSupply struct {
	ctx context.Context
}

func TokenSet() {

}

func (t *TotalSupply) Start() {

}
