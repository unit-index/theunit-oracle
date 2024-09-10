package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/toknowwhy/theunit-oracle/internal/gofer/marshal"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	pkgEthereum "github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
	"github.com/toknowwhy/theunit-oracle/pkg/oracle"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/messages"
	"github.com/toknowwhy/theunit-oracle/pkg/unit"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/feeder"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/nodes"
	"math/big"
	"sync"
	"time"
)

var errInvalidSignature = errors.New("received unitParams has an invalid signature")
var errUnknownToken = errors.New("received token is not configured")
var errUnknownFeeder = errors.New("feeder is not allowed to send unitParams")
var errInvalidParams = errors.New("received params is invalid")

type ErrTokenNotFound struct {
	Token unit.Token
}

func (e ErrTokenNotFound) Error() string {
	return fmt.Sprintf("unable to find the %s pair", e.Token)
}

type FeederUnitParams struct {
	TokenName string
	Feeder    ethereum.Address
}

type UnitParamStore struct {
	mu         sync.RWMutex
	unitParams map[FeederUnitParams]*messages.UnitParams
}

func NewUnitParamStore() *UnitParamStore {
	return &UnitParamStore{
		unitParams: make(map[FeederUnitParams]*messages.UnitParams),
	}
}

func (u *UnitParamStore) Add(from ethereum.Address, msg *messages.UnitParams) {
	u.mu.Lock()
	defer u.mu.Unlock()

	fp := FeederUnitParams{
		TokenName: msg.UnitParams.Name,
		Feeder:    from,
	}

	if prev, ok := u.unitParams[fp]; ok && prev.UnitParams.Age.After(msg.UnitParams.Age) {
		return
	}

	u.unitParams[fp] = msg
}

type Unit struct {
	graphs         map[unit.Token]nodes.Aggregator
	feeder         *feeder.Feeder
	unitAlgorithm  oracle.UnitAlgorithm
	gofer          gofer.Gofer
	tokens         map[string]unit.Token
	doneCh         chan struct{}
	interval       time.Duration
	signer         ethereum.Signer
	transport      transport.Transport
	mu             sync.Mutex
	ctx            context.Context
	Feeds          []ethereum.Address
	UnitParamStore UnitParamStore
}

func NewUnit(
	g map[unit.Token]nodes.Aggregator,
	f *feeder.Feeder,
	unitAlgorithm oracle.UnitAlgorithm,
	gofer gofer.Gofer,
	tokens map[string]unit.Token,
	interval time.Duration,
	signer pkgEthereum.Signer,
	transport transport.Transport,
	feedAddresses []ethereum.Address,
	unitParamStore UnitParamStore,
) *Unit {
	ctx := context.Background()
	return &Unit{
		graphs:         g,
		feeder:         f,
		unitAlgorithm:  unitAlgorithm,
		gofer:          gofer,
		tokens:         tokens,
		doneCh:         make(chan struct{}),
		interval:       interval,
		signer:         signer,
		transport:      transport,
		ctx:            ctx,
		Feeds:          feedAddresses,
		UnitParamStore: unitParamStore,
	}
}

func (u *Unit) TokenTotalSupply(token unit.Token) (*unit.CSupply, error) {
	n, ok := u.graphs[token]
	if !ok {
		return nil, ErrTokenNotFound{Token: token}
	}

	if u.feeder != nil {
		u.feeder.Feed(n)
	}
	return mapGraphCSupply(n.CSupply()), nil
}

func (u *Unit) TokensTotalSupply(tokens ...unit.Token) (map[unit.Token]*unit.CSupply, error) {
	ns, err := u.findNodes(tokens...)

	if err != nil {
		return nil, err
	}
	if u.feeder != nil {
		errs := u.feeder.Feed(ns...)
		if len(errs.List) != 0 {
			fmt.Println(errs)
		}
	}

	res := make(map[unit.Token]*unit.CSupply)
	for _, n := range ns {
		if n, ok := n.(nodes.Aggregator); ok {
			res[n.Token()] = mapGraphCSupply(n.CSupply())
		}
	}
	return res, nil
}

func (u *Unit) FeedMarketCapAndPrice(tokens ...unit.Token) (map[string]unit.UnitPerMonthParams, error) {
	css, err := u.TokensTotalSupply(tokens...)
	if err != nil {
		return nil, err
	}

	//var unitPerMonthParams []unit.UnitPerMonthParams
	var unitPerMonthParams = make(map[string]unit.UnitPerMonthParams)
	//  去合约里查本月有哪些token在index内
	//获取上个月最后一天的时间戳。然后再去合约里查
	now := time.Now()
	firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfLastMonth := firstDayOfThisMonth.Add(-time.Second)
	lastDayOfLastMonthUTC := lastDayOfLastMonth.In(time.UTC)

	timestamp := lastDayOfLastMonthUTC.Unix()
	bigIntTimestamp := big.NewInt(0).SetInt64(timestamp)
	ctx := context.Background()
	token, err := u.unitAlgorithm.GetTokens(ctx, bigIntTimestamp)
	if err != nil {
		return nil, err
	}
	fmt.Println(token)
	// 检测下配置里的token与合约里的，只能比合约里的多，不能少于合约。
	// 去调合约
	for _, token := range u.tokens {
		//csupply, err := u.TokenTotalSupply(token)
		//fmt.Println("marketCap", token.Symbol)
		if err != nil {
			return nil, err
		}
		pair, err := gofer.NewPair(fmt.Sprintf("%s/%s", token.Symbol, "USD"))
		if err != nil {
			return nil, err
		}
		p, err := u.gofer.Price(pair)
		if err != nil {
			return nil, err
		}
		CSupply := css[token].CSupply
		marketCap := CSupply * p.Price

		unitPerMonthParams[token.Name] = unit.UnitPerMonthParams{CSupply: CSupply, LastPrice: p.Price, LastMarketCap: marketCap}
	}

	return unitPerMonthParams, nil
}

func (u *Unit) Price() (string, error) {
	//now := time.Now()
	//firstDayOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	//lastDayOfLastMonth := firstDayOfThisMonth.Add(-time.Second)
	//lastDayOfLastMonthUTC := lastDayOfLastMonth.In(time.UTC)
	//timestamp := lastDayOfLastMonthUTC.Unix()

	return "", nil
}

func (u *Unit) Start() error {
	u.broadcasterLoop()
	u.collectorLoop()
	u.relayerLoop()
	return nil
}

func (u *Unit) collectorLoop() error {
	go func() {
		u.mu.Lock()
		defer u.mu.Unlock()
		for {
			select {
			case <-u.ctx.Done():
				return
			case m := <-u.transport.Messages(messages.UnitParamsMessageName):
				// If there was a problem while reading prices from the transport:
				if m.Error != nil {
					//u.log.WithError(m.Error).Warn("Unable to read prices from the transport")
					fmt.Println(m.Error)
					continue
				}
				unitParams, ok := m.Message.(*messages.UnitParams)
				if !ok {
					fmt.Println(ok)
					//u.log.Error("Unexpected value returned from transport layer")
					continue
				}

				// Try to collect received price:
				err := u.collectUnitParams(unitParams)
				// Print logs:
				if err != nil {
					//u.log.
					//	WithError(err).
					//	WithFields(price.Price.Fields(c.signer)).
					//	Warn("Received invalid price")
				} else {
					//u.log.
					//	WithFields(price.Price.Fields(c.signer)).
					//	Info("Price received")
				}
			}
		}
	}()

	return nil
}

func (u *Unit) collectUnitParams(msg *messages.UnitParams) error {
	from, err := msg.UnitParams.From(u.signer)

	if err != nil {
		return errInvalidSignature
	}

	if _, ok := u.tokens[msg.UnitParams.Name]; !ok {
		return errUnknownToken
	}
	if !u.isFeedAllowed(*from) {
		return errUnknownFeeder
	}
	if msg.UnitParams.LastPrice.Cmp(big.NewInt(0)) <= 0 || msg.UnitParams.LastMarketCap.Cmp(big.NewInt(0)) <= 0 {
		return errInvalidParams
	}
	fmt.Println("collectUnitParams", from, msg.UnitParams.Age, msg.UnitParams.LastMarketCap, msg.UnitParams.Name, msg.UnitParams.LastPrice)
	u.UnitParamStore.Add(*from, msg)

	return nil
}
func (u *Unit) isFeedAllowed(address ethereum.Address) bool {
	for _, a := range u.Feeds {
		if a == address {
			return true
		}
	}
	return false
}

func (u *Unit) broadcasterLoop() error {
	if u.interval == 0 {
		return nil
	}
	ticker := time.NewTicker(u.interval)
	wg := sync.WaitGroup{}
	go func() {
		for {
			select {
			case <-u.doneCh:
				ticker.Stop()
				return
			case <-ticker.C:
				wg.Add(1)
				go func() {

					unitPerMonthParams, err := u.FeedMarketCapAndPrice()
					if err != nil {
						fmt.Println(err)
					}
					err = u.broadcast(unitPerMonthParams)
					if err != nil {
						fmt.Println(err)
					}
					//for assetPair := range g.goferPairs {
					//	err := g.broadcast(assetPair)
					//	if err != nil {
					//		g.log.
					//			WithFields(log.Fields{"assetPair": assetPair}).
					//			WithError(err).
					//			Warn("Unable to broadcast price")
					//	} else {
					//		g.log.
					//			WithFields(log.Fields{"assetPair": assetPair}).
					//			Info("Price broadcast")
					//	}
					//}
					wg.Done()
				}()
			}

			wg.Wait()
		}
	}()
	return nil
}

func (u *Unit) broadcast(unitPerMonthParams map[string]unit.UnitPerMonthParams) error {
	fmt.Println("broadcast successful!", u.tokens)
	for _, token := range u.tokens {
		upmp, _ := unitPerMonthParams[token.Name]

		LastMarketCap := new(big.Int)
		int64LastMarketCap := int64(upmp.LastMarketCap)
		LastMarketCap.SetInt64(int64LastMarketCap)

		LastPrice := new(big.Int)
		int64LastPrice := int64(upmp.LastPrice)
		LastPrice.SetInt64(int64LastPrice)

		unitParams := &oracle.UnitParams{Name: token.Name, LastMarketCap: LastMarketCap, LastPrice: LastPrice, Age: time.Now()}
		fmt.Println("broadcast successful0000!")
		err := unitParams.Sign(u.signer)

		if err != nil {
			return err
		}

		message, err := createUnitParamsMessage(unitParams, &upmp)
		fmt.Println("broadcast successful22222!", err)
		if err != nil {
			return err
		}

		err = u.transport.Broadcast(messages.UnitParamsMessageName, message)
		fmt.Println("broadcast successful!", err)
		if err != nil {
			return err
		}
	}

	return nil
}

func createUnitParamsMessage(up *oracle.UnitParams, upmp *unit.UnitPerMonthParams) (*messages.UnitParams, error) {
	trace, err := marshal.Marshall(marshal.JSON, upmp)
	if err != nil {
		return nil, err
	}

	return &messages.UnitParams{
		UnitParams: up,
		Trace:      trace,
	}, nil
}

func (u *Unit) relayerLoop() error {
	if u.interval == 0 {
		return nil
	}
	ticker := time.NewTicker(u.interval)

	wg := sync.WaitGroup{}
	go func() {
		for {
			select {
			case <-u.doneCh:
				ticker.Stop()
				return
			case <-ticker.C:
				// TODO: fetch all prices before broadcast is called
				wg.Add(1)
				go func() {
					//unitPerMonthParams, err := u.FeedMarketCapAndPrice()
					//if err != nil {
					//	u.relayer(unitPerMonthParams)
					//}

					wg.Done()
				}()
			}
			wg.Wait()
		}
	}()
}

//func (u *Unit) relayer(unitPerMonthParams []unit.UnitPerMonthParams) error {
//
//}

func (u *Unit) Wait() {}

func mapGraphCSupply(t interface{}) *unit.CSupply {
	gt := &unit.CSupply{
		Parameters: make(map[string]string),
	}

	switch typedCSupply := t.(type) {
	case nodes.AggregatorCSupply:
		gt.Type = "aggregator"
		gt.Token = typedCSupply.Token
		gt.CSupply = typedCSupply.CSupply
		gt.Time = typedCSupply.Time
		if typedCSupply.Error != nil {
			gt.Error = typedCSupply.Error.Error()
		}
		gt.Parameters = typedCSupply.Parameters
		for _, ct := range typedCSupply.OriginCSupply {
			gt.CSupplys = append(gt.CSupplys, mapGraphCSupply(ct))
		}
		for _, ct := range typedCSupply.AggregatorPrices {
			gt.CSupplys = append(gt.CSupplys, mapGraphCSupply(ct))
		}
	case nodes.OriginCSupply:
		gt.Type = "origin"
		gt.Token = typedCSupply.Token
		gt.CSupply = typedCSupply.CSupply
		gt.Time = typedCSupply.Time
		if typedCSupply.Error != nil {
			gt.Error = typedCSupply.Error.Error()
		}
		gt.Parameters["origin"] = typedCSupply.Origin
	default:
		panic("unsupported object")
	}

	return gt
}

func (u *Unit) findNodes(tokens ...unit.Token) ([]nodes.Node, error) {
	var ns []nodes.Node
	if len(tokens) == 0 { // Return all:
		for _, n := range u.graphs {
			ns = append(ns, n)
		}
	} else { // Return for given pairs:
		for _, p := range tokens {
			n, ok := u.graphs[p]
			if !ok {
				return nil, ErrTokenNotFound{Token: p}
			}
			ns = append(ns, n)
		}
	}
	return ns, nil
}
