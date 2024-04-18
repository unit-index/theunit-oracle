package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/pkg/datastore/memory/testutil"
)

func TestPriceStore_Add(t *testing.T) {
	ps := NewPriceStore()

	ps.Add(testutil.Address1, testutil.PriceAAABBB1)
	ps.Add(testutil.Address1, testutil.PriceXXXYYY1)
	ps.Add(testutil.Address2, testutil.PriceAAABBB1)
	ps.Add(testutil.Address2, testutil.PriceXXXYYY1)

	aaabbb := ps.AssetPair("AAABBB")
	xxxyyy := ps.AssetPair("XXXYYY")

	assert.Equal(t, 2, len(aaabbb))
	assert.Equal(t, 2, len(xxxyyy))
	assert.Contains(t, aaabbb, testutil.PriceAAABBB1)
	assert.Contains(t, xxxyyy, testutil.PriceXXXYYY1)
}

func TestPriceStore_Add_UseNewerPrice(t *testing.T) {
	ps := NewPriceStore()

	// Second price should replace first one because is younger:
	ps.Add(testutil.Address1, testutil.PriceAAABBB1)
	ps.Add(testutil.Address1, testutil.PriceAAABBB2)

	// Second price should be ignored because is older:
	ps.Add(testutil.Address1, testutil.PriceXXXYYY2)
	ps.Add(testutil.Address1, testutil.PriceXXXYYY1)

	aaabbb := ps.AssetPair("AAABBB")
	xxxyyy := ps.AssetPair("XXXYYY")

	assert.Equal(t, testutil.PriceAAABBB2, aaabbb[0])
	assert.Equal(t, testutil.PriceXXXYYY2, xxxyyy[0])
}

func TestPriceStore_Feeder(t *testing.T) {
	ps := NewPriceStore()

	ps.Add(testutil.Address1, testutil.PriceAAABBB1)
	ps.Add(testutil.Address1, testutil.PriceAAABBB2)
	ps.Add(testutil.Address1, testutil.PriceXXXYYY1)
	ps.Add(testutil.Address1, testutil.PriceXXXYYY2)
	ps.Add(testutil.Address2, testutil.PriceAAABBB1)
	ps.Add(testutil.Address2, testutil.PriceAAABBB2)
	ps.Add(testutil.Address2, testutil.PriceXXXYYY1)
	ps.Add(testutil.Address2, testutil.PriceXXXYYY2)

	assert.Equal(t, testutil.PriceAAABBB2, ps.Feeder("AAABBB", testutil.Address1))
	assert.Equal(t, testutil.PriceXXXYYY2, ps.Feeder("XXXYYY", testutil.Address1))
}
