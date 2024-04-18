package origins

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	ethereumMocks "github.com/toknowwhy/theunit-oracle/pkg/ethereum/mocks"

	"github.com/stretchr/testify/suite"
)

type BalancerV2Suite struct {
	suite.Suite
	addresses ContractAddresses
	client    *ethereumMocks.Client
	origin    *BaseExchangeHandler
}

func (suite *BalancerV2Suite) SetupSuite() {
	suite.addresses = ContractAddresses{
		"STETH/ETH": "0x32296969ef14eb0c6d29669c550d4a0449130230",
	}
	suite.client = &ethereumMocks.Client{}
}
func (suite *BalancerV2Suite) TearDownSuite() {
	suite.addresses = nil
	suite.client = nil
}

func (suite *BalancerV2Suite) SetupTest() {
	balancerV2Finance, err := NewBalancerV2(suite.client, suite.addresses)
	suite.NoError(err)
	suite.origin = NewBaseExchangeHandler(balancerV2Finance, nil)
}

func (suite *BalancerV2Suite) TearDownTest() {
	suite.origin = nil
}

func (suite *BalancerV2Suite) Origin() Handler {
	return suite.origin
}

func TestBalancerV2Suite(t *testing.T) {
	suite.Run(t, new(BalancerV2Suite))
}

func (suite *BalancerV2Suite) TestSuccessResponse() {
	suite.client.On("Call", mock.Anything, ethereum.Call{
		Address: ethereum.HexToAddress("0x32296969Ef14EB0c6d29669C550D4a0449130230"),
		Data:    ethereum.HexToBytes("0xb10be7390000000000000000000000000000000000000000000000000000000000000000"),
	}).Return(ethereum.HexToBytes("0x0000000000000000000000000000000000000000000000000dc19f91822f3fe3"), nil)

	pair := Pair{Base: "STETH", Quote: "ETH"}

	results1 := suite.origin.Fetch([]Pair{pair})
	suite.Require().NoError(results1[0].Error)
	suite.Equal(0.9912488403014287, results1[0].Price.Price)
	suite.Greater(results1[0].Price.Timestamp.Unix(), int64(0))

	results2 := suite.origin.Fetch([]Pair{pair.Inverse()})
	suite.Require().Error(results2[0].Error)
}

func (suite *BalancerV2Suite) TestFailOnWrongPair() {
	pair := Pair{Base: "x", Quote: "y"}
	cr := suite.origin.Fetch([]Pair{pair})
	suite.Require().EqualError(cr[0].Error, "failed to get contract address for pair: x/y")
}
