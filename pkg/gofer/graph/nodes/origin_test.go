package nodes

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

const originTestTTL = 10 * time.Second

func TestOriginNode_OriginPair(t *testing.T) {
	op := OriginPair{
		Origin: "foo",
		Pair:   gofer.Pair{Base: "A", Quote: "B"},
	}

	o := NewOriginNode(op, originTestTTL, originTestTTL)
	assert.Equal(t, op, o.OriginPair())
}

func TestOriginNode_Ingest_Valid(t *testing.T) {
	op := OriginPair{
		Origin: "foo",
		Pair:   gofer.Pair{Base: "A", Quote: "B"},
	}

	ot := OriginPrice{
		PairPrice: PairPrice{
			Pair:      gofer.Pair{Base: "A", Quote: "B"},
			Price:     10,
			Bid:       10,
			Ask:       10,
			Volume24h: 10,
			Time:      time.Now(),
		},
		Origin: "foo",
		Error:  nil,
	}

	o := NewOriginNode(op, originTestTTL, originTestTTL)
	err := o.Ingest(ot)

	assert.Equal(t, op, o.OriginPair())
	assert.Equal(t, ot, o.Price())
	assert.NoError(t, err)
	assert.NoError(t, o.price.Error)
}

func TestOriginNode_Ingest_IncompatiblePair(t *testing.T) {
	op := OriginPair{
		Origin: "foo",
		Pair:   gofer.Pair{Base: "A", Quote: "B"},
	}

	ot := OriginPrice{
		PairPrice: PairPrice{
			Pair:      gofer.Pair{Base: "A", Quote: "C"},
			Price:     10,
			Bid:       10,
			Ask:       10,
			Volume24h: 10,
			Time:      time.Now(),
		},
		Origin: "foo",
		Error:  nil,
	}

	o := NewOriginNode(op, originTestTTL, originTestTTL)
	err := o.Ingest(ot)

	assert.True(t, errors.As(err, &ErrIncompatiblePair{}))
	assert.NoError(t, o.price.Error)
}

func TestOriginNode_Ingest_IncompatibleOrigin(t *testing.T) {
	op := OriginPair{
		Origin: "foo",
		Pair:   gofer.Pair{Base: "A", Quote: "B"},
	}

	ot := OriginPrice{
		PairPrice: PairPrice{
			Pair:      gofer.Pair{Base: "A", Quote: "B"},
			Price:     10,
			Bid:       10,
			Ask:       10,
			Volume24h: 10,
			Time:      time.Now(),
		},
		Origin: "bar",
		Error:  nil,
	}

	o := NewOriginNode(op, originTestTTL, originTestTTL)
	err := o.Ingest(ot)

	assert.True(t, errors.As(err, &IncompatibleOriginErr{}))
	assert.NoError(t, o.price.Error)
}

func TestOriginNode_Ingest_IncompatibleEverything(t *testing.T) {
	op := OriginPair{
		Origin: "foo",
		Pair:   gofer.Pair{Base: "A", Quote: "B"},
	}

	ot := OriginPrice{
		PairPrice: PairPrice{
			Pair:      gofer.Pair{Base: "A", Quote: "C"},
			Price:     10,
			Bid:       10,
			Ask:       10,
			Volume24h: 10,
			Time:      time.Now(),
		},
		Origin: "bar",
		Error:  nil,
	}

	o := NewOriginNode(op, originTestTTL, originTestTTL)
	err := o.Ingest(ot)

	assert.True(t, errors.As(err, &IncompatibleOriginErr{}))
	assert.True(t, errors.As(err, &ErrIncompatiblePair{}))
	assert.NoError(t, o.price.Error)
}

func TestOriginNode_Ingest_PriceWithError(t *testing.T) {
	err := errors.New("something")

	op := OriginPair{
		Origin: "foo",
		Pair:   gofer.Pair{Base: "A", Quote: "B"},
	}

	ot := OriginPrice{
		PairPrice: PairPrice{
			Pair:      gofer.Pair{Base: "A", Quote: "B"},
			Price:     10,
			Bid:       10,
			Ask:       10,
			Volume24h: 10,
			Time:      time.Now(),
		},
		Origin: "foo",
		Error:  err,
	}

	o := NewOriginNode(op, originTestTTL, originTestTTL)
	err2 := o.Ingest(ot)

	assert.NoError(t, err2)
	assert.Equal(t, err, o.price.Error)
}

func TestOriginNode_Price_Expired(t *testing.T) {
	op := OriginPair{
		Origin: "foo",
		Pair:   gofer.Pair{Base: "A", Quote: "B"},
	}

	ot := OriginPrice{
		PairPrice: PairPrice{
			Pair:      gofer.Pair{Base: "A", Quote: "B"},
			Price:     10,
			Bid:       10,
			Ask:       10,
			Volume24h: 10,
			Time:      time.Now().Add(-20 * time.Second),
		},
		Origin: "foo",
		Error:  nil,
	}

	o := NewOriginNode(op, originTestTTL, originTestTTL)
	_ = o.Ingest(ot)
	price := o.Price()

	assert.True(t, errors.As(price.Error, &ErrPriceTTLExpired{}))
}
