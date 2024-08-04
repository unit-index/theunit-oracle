package feeder

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/toknowwhy/theunit-oracle/pkg/gofer/graph/nodes"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/unit/origins"
	"time"
)

const LoggerTag = "FEEDER"

type Warnings struct {
	List []error
}

func (w Warnings) ToError() error {
	var err error
	for _, e := range w.List {
		err = multierror.Append(err, e)
	}
	return err
}

type Feedable interface {
	// OriginPair returns the origin and pair which are acceptable for
	// this Node.
	OriginPair() nodes.OriginPair
	// Ingest sets the Price for this Node. It may return error if
	// the OriginPrice contains incompatible origin or pair.
	Ingest(price nodes.OriginPrice) error
	// MinTTL is the amount of time during which the Price shouldn't be updated.
	MinTTL() time.Duration
	// MaxTTL is the maximum amount of time during which the Price can be used.
	// After that time, the Price method will return a OriginPrice with
	// a ErrPriceTTLExpired error.
	MaxTTL() time.Duration
	// Expired returns true if the Price is expired. This is based on the MaxTTL
	// value.
	Expired() bool
	// Price returns the Price assigned in the Ingest method. If the Price is
	// expired then a ErrPriceTTLExpired error will be set in
	// the OriginPrice.Error field.
	Price() nodes.OriginPrice
}

type Feeder struct {
	ctx context.Context

	set    *origins.Set
	log    log.Logger
	doneCh chan struct{}
}

// NewFeeder creates new Feeder instance.
func NewFeeder(ctx context.Context, set *origins.Set, log log.Logger) *Feeder {
	return &Feeder{
		ctx:    ctx,
		set:    set,
		log:    log.WithField("tag", LoggerTag),
		doneCh: make(chan struct{}),
	}
}

func (f *Feeder) Start(ns ...nodes.Node) error {
	f.log.Infof("Starting")

	gcdTTL := getGCDTTL(ns)
	if gcdTTL < time.Second {
		gcdTTL = time.Second
	}
	f.log.WithField("interval", gcdTTL.String()).Infof("Update interval (GCD of all TTLs)")

	feed := func() {
		// We have to add gcdTTL to the current time because we want
		// to find all nodes that will expire before the next tick.
		t := time.Now().Add(gcdTTL)
		warns := f.fetchSupplyAndPriceThemToFeedableNodes(f.findFeedableNodes(ns, t))
		if len(warns.List) > 0 {
			f.log.WithError(warns.ToError()).Warn("Unable to feed some nodes")
		}
	}

	ticker := time.NewTicker(gcdTTL)
	go func() {
		feed()
		for {
			select {
			case <-f.ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				feed()
			}
		}
	}()

	go f.contextCancelHandler()
	return nil
}
func (f *Feeder) Wait() {
	<-f.doneCh
}

func (f *Feeder) contextCancelHandler() {
	defer func() { close(f.doneCh) }()
	defer f.log.Info("Stopped")

	<-f.ctx.Done()
}
func (f *Feeder) findFeedableNodes(ns []nodes.Node, t time.Time) []Feedable {
	var feedables []Feedable
	nodes.Walk(func(n nodes.Node) {
		if feedable, ok := n.(Feedable); ok {
			if t.Sub(feedable.Price().Time) >= feedable.MinTTL() {
				feedables = append(feedables, feedable)
			}
		}
	}, ns...)

	return feedables
}

//func (f *Feeder) Feed(ns ...nodes.Node) Warnings {
//	return f.fetchSupplyAndPriceThemToFeedableNodes(ns)
//}

func (f *Feeder) fetchSupplyAndPriceThemToFeedableNodes(ns []Feedable) Warnings {
	var warn Warnings

	return warn
}

func getGCDTTL(ns []nodes.Node) time.Duration {
	ttl := time.Duration(0)
	nodes.Walk(func(n nodes.Node) {
		if f, ok := n.(Feedable); ok {
			if ttl == 0 {
				ttl = f.MinTTL()
			}
			a := ttl
			b := f.MinTTL()
			for b != 0 {
				t := b
				b = a % b
				a = t
			}
			ttl = a
		}
	}, ns...)
	return ttl
}
