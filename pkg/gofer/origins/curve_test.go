package origins

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	ethereumMocks "github.com/toknowwhy/theunit-oracle/pkg/ethereum/mocks"

	"github.com/stretchr/testify/suite"
)

type CurveSuite struct {
	suite.Suite
	addresses ContractAddresses
	client    *ethereumMocks.Client
	origin    *BaseExchangeHandler
}

func (suite *CurveSuite) SetupSuite() {
	suite.addresses = ContractAddresses{
		"ETH/STETH": "0xDC24316b9AE028F1497c275EB9192a3Ea0f67022",
	}
	suite.client = &ethereumMocks.Client{}
}
func (suite *CurveSuite) TearDownSuite() {
	suite.addresses = nil
	suite.client = nil
}

func (suite *CurveSuite) SetupTest() {
	curveFinance, err := NewCurveFinance(suite.client, suite.addresses)
	suite.NoError(err)
	suite.origin = NewBaseExchangeHandler(curveFinance, nil)
}

func (suite *CurveSuite) TearDownTest() {
	suite.origin = nil
}

func (suite *CurveSuite) Origin() Handler {
	return suite.origin
}

func TestCurveSuite(t *testing.T) {
	suite.Run(t, new(CurveSuite))
}

func (suite *CurveSuite) TestSuccessResponse() {
	suite.client.On("Call", mock.Anything, ethereum.Call{
		Address: ethereum.HexToAddress("0xDC24316b9AE028F1497c275EB9192a3Ea0f67022"),
		Data:    ethereum.HexToBytes("0x5e0d443f000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000de0b6b3a7640000"),
	}).Return(ethereum.HexToBytes("0x0000000000000000000000000000000000000000000000000dc19f91822f3fe3"), nil)

	pair := Pair{Base: "STETH", Quote: "ETH"}

	results1 := suite.origin.Fetch([]Pair{pair})
	suite.Require().NoError(results1[0].Error)
	suite.Equal(0.9912488403014287, results1[0].Price.Price)
	suite.Greater(results1[0].Price.Timestamp.Unix(), int64(0))

	suite.client.On("Call", mock.Anything, ethereum.Call{
		Address: ethereum.HexToAddress("0xDC24316b9AE028F1497c275EB9192a3Ea0f67022"),
		Data:    ethereum.HexToBytes("0x5e0d443f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000de0b6b3a7640000"),
	}).Return(ethereum.HexToBytes("0x0000000000000000000000000000000000000000000000000dc19f91822f3fe3"), nil)

	results2 := suite.origin.Fetch([]Pair{pair.Inverse()})
	suite.Require().NoError(results2[0].Error)
	suite.Equal(0.9912488403014287, results2[0].Price.Price)
	suite.Greater(results2[0].Price.Timestamp.Unix(), int64(0))
}

func (suite *CurveSuite) TestFailOnWrongPair() {
	pair := Pair{Base: "x", Quote: "y"}
	cr := suite.origin.Fetch([]Pair{pair})
	suite.Require().EqualError(cr[0].Error, "failed to get contract address for pair: x/y")
}
