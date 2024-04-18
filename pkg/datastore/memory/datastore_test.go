package memory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/toknowwhy/theunit-oracle/pkg/datastore/memory/testutil"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum/mocks"
	"github.com/toknowwhy/theunit-oracle/pkg/log/null"
	"github.com/toknowwhy/theunit-oracle/pkg/oracle"
	"github.com/toknowwhy/theunit-oracle/pkg/transport"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/local"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/messages"
)

func TestDatastore_Prices(t *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	sig := &mocks.Signer{}
	tra := local.New(ctx, 0, map[string]transport.Message{messages.PriceMessageName: (*messages.Price)(nil)})

	ds, err := NewDatastore(ctx, Config{
		Signer:    sig,
		Transport: tra,
		Pairs: map[string]*Pair{
			"AAABBB": {Feeds: []ethereum.Address{testutil.Address1, testutil.Address2}},
			"XXXYYY": {Feeds: []ethereum.Address{testutil.Address1, testutil.Address2}},
		},
		Logger: null.New(),
	})
	require.NoError(t, err)
	require.NoError(t, ds.Start())

	sig.On("Recover", testutil.PriceAAABBB1.Price.Signature(), mock.Anything).Return(&testutil.Address1, nil)
	sig.On("Recover", testutil.PriceAAABBB2.Price.Signature(), mock.Anything).Return(&testutil.Address2, nil)
	sig.On("Recover", testutil.PriceXXXYYY1.Price.Signature(), mock.Anything).Return(&testutil.Address1, nil)
	sig.On("Recover", testutil.PriceXXXYYY2.Price.Signature(), mock.Anything).Return(&testutil.Address2, nil)

	assert.NoError(t, tra.Broadcast(messages.PriceMessageName, testutil.PriceAAABBB1))
	assert.NoError(t, tra.Broadcast(messages.PriceMessageName, testutil.PriceAAABBB2))
	assert.NoError(t, tra.Broadcast(messages.PriceMessageName, testutil.PriceXXXYYY1))
	assert.NoError(t, tra.Broadcast(messages.PriceMessageName, testutil.PriceXXXYYY2))

	// Datastore fetches prices asynchronously, so we wait up to 1 second:
	var aaabbb, xxxyyy []*messages.Price
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		aaabbb = ds.Prices().AssetPair("AAABBB")
		xxxyyy = ds.Prices().AssetPair("XXXYYY")
		if len(aaabbb) == 2 && len(xxxyyy) == 2 {
			break
		}
	}

	assert.Contains(t, toOraclePrices(aaabbb), testutil.PriceAAABBB1.Price)
	assert.Contains(t, toOraclePrices(aaabbb), testutil.PriceAAABBB2.Price)
	assert.Contains(t, toOraclePrices(xxxyyy), testutil.PriceXXXYYY1.Price)
	assert.Contains(t, toOraclePrices(xxxyyy), testutil.PriceXXXYYY2.Price)
}

func toOraclePrices(ps []*messages.Price) []*oracle.Price {
	var r []*oracle.Price
	for _, p := range ps {
		r = append(r, p.Price)
	}
	return r
}
