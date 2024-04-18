package memory

import (
	"context"
	"errors"
	"math/big"
	"sync"

	"github.com/toknowwhy/theunit-oracle/pkg/datastore"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/messages"
)

const LoggerTag = "DATASTORE"

var errInvalidSignature = errors.New("received price has an invalid signature")
var errInvalidPrice = errors.New("received price is invalid")
var errUnknownPair = errors.New("received pair is not configured")
var errUnknownFeeder = errors.New("feeder is not allowed to send prices")

// Datastore reads and stores prices from the P2P network.
type Datastore struct {
	ctx    context.Context
	mu     sync.Mutex
	doneCh chan struct{}

	signer     ethereum.Signer
	transport  transport.Transport
	pairs      map[string]*Pair
	priceStore *PriceStore
	log        log.Logger
}

type Config struct {
	// Signer is an instance of the ethereum.Signer which will be used to
	// verify price signatures.
	Signer ethereum.Signer
	// Transport is a implementation of transport used to fetch prices from
	// feeders.
	Transport transport.Transport
	// Pairs is the list supported pairs by the datastore with their
	// configuration.
	Pairs map[string]*Pair
	// Logger is a current logger interface used by the Datastore.
	// The Logger is required to monitor asynchronous processes.
	Logger log.Logger
}

type Pair struct {
	// Feeds is the list of Ethereum addresses from which prices will be
	// accepted.
	Feeds []ethereum.Address
}

func NewDatastore(ctx context.Context, cfg Config) (*Datastore, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}
	return &Datastore{
		ctx:        ctx,
		doneCh:     make(chan struct{}),
		signer:     cfg.Signer,
		transport:  cfg.Transport,
		pairs:      cfg.Pairs,
		priceStore: NewPriceStore(),
		log:        cfg.Logger.WithField("tag", LoggerTag),
	}, nil
}

// Start implements the datastore.Datastore interface.
func (c *Datastore) Start() error {
	c.log.Info("Starting")

	go c.contextCancelHandler()
	return c.collectorLoop()
}

// Wait implements the datastore.Datastore interface.
func (c *Datastore) Wait() {
	<-c.doneCh
}

// Prices implements the datastore.Datastore interface.
func (c *Datastore) Prices() datastore.PriceStore {
	return c.priceStore
}

// collectPrice adds a price from a feeder which may be used to update
// Oracle contract. The price will be added only if a feeder is
// allowed to send prices.
func (c *Datastore) collectPrice(msg *messages.Price) error {
	from, err := msg.Price.From(c.signer)
	if err != nil {
		return errInvalidSignature
	}
	if _, ok := c.pairs[msg.Price.Wat]; !ok {
		return errUnknownPair
	}
	if !c.isFeedAllowed(msg.Price.Wat, *from) {
		return errUnknownFeeder
	}
	if msg.Price.Val.Cmp(big.NewInt(0)) <= 0 {
		return errInvalidPrice
	}

	c.priceStore.Add(*from, msg)

	return nil
}

// collectorLoop creates a asynchronous loop which fetches prices from feeders.
func (c *Datastore) collectorLoop() error {
	go func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		for {
			select {
			case <-c.ctx.Done():
				return
			case m := <-c.transport.Messages(messages.PriceMessageName):
				// If there was a problem while reading prices from the transport:
				if m.Error != nil {
					c.log.
						WithError(m.Error).
						Warn("Unable to read prices from the transport")
					continue
				}
				price, ok := m.Message.(*messages.Price)
				if !ok {
					c.log.Error("Unexpected value returned from transport layer")
					continue
				}

				// Try to collect received price:
				err := c.collectPrice(price)

				// Print logs:
				if err != nil {
					c.log.
						WithError(err).
						WithFields(price.Price.Fields(c.signer)).
						Warn("Received invalid price")
				} else {
					c.log.
						WithFields(price.Price.Fields(c.signer)).
						Info("Price received")
				}
			}
		}
	}()

	return nil
}

func (c *Datastore) isFeedAllowed(assetPair string, address ethereum.Address) bool {
	for _, a := range c.pairs[assetPair].Feeds {
		if a == address {
			return true
		}
	}
	return false
}

func (c *Datastore) contextCancelHandler() {
	defer func() { close(c.doneCh) }()
	defer c.log.Info("Stopped")
	<-c.ctx.Done()
}
