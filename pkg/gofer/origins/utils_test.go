package origins

import (
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite interface {
	suite.TestingSuite

	Assert() *assert.Assertions
	Origin() Handler
}

func testRealAPICall(suite Suite, origin *BaseExchangeHandler, base, quote string) {
	testRealBatchAPICall(suite, origin, []Pair{{Base: base, Quote: quote}})
}

func testRealBatchAPICall(suite Suite, origin *BaseExchangeHandler, pairs []Pair) {
	if os.Getenv("GOFER_TEST_API_CALLS") == "" {
		suite.T().SkipNow()
	}

	suite.Assert().IsType(suite.Origin(), origin)

	crs := origin.Fetch(pairs)

	for _, cr := range crs {
		suite.Assert().NoErrorf(cr.Error, "%q", cr.Price.Pair)
		suite.Assert().Greater(cr.Price.Price, float64(0))
	}
}
