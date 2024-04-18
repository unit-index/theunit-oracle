package oracle

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

// Median is an interface for the median oracle contract:
// https://github.com/makerdao/median/
//
// Contract documentation:
// https://docs.makerdao.com/smart-contract-modules/oracle-module/median-detailed-documentation
type Median interface {
	// Address returns medianizer contract address.
	Address() common.Address
	// Age returns the value from contract's age method. The age is the block
	// timestamp of last price val update.
	Age(ctx context.Context) (time.Time, error)
	// Bar returns the value from contract's bar method. The bar method returns
	// the minimum number of prices necessary to accept a new median value.
	Bar(ctx context.Context) (int64, error)
	// Val returns current asset price form the contract's storage.
	Val(ctx context.Context) (*big.Int, error)
	// Wat returns asset name.
	Wat(ctx context.Context) (string, error)
	// Feeds returns a list of all Ethereum addresses that are authorized to update
	// Oracle prices (orcls).
	Feeds(ctx context.Context) ([]ethereum.Address, error)
	// Poke sends transaction to the smart contract which invokes contract's
	// poke method, which updates asset price (val).  If simulateBeforeRun is
	// set to true, then transaction will be simulated on the EVM before actual
	// transaction will be send.
	Poke(ctx context.Context, prices []*Price, simulateBeforeRun bool) (*ethereum.Hash, error)
	// Lift sends transaction to the smart contract which invokes contract's
	// lift method, which sends  adds given addresses to the feeders list (orcls).
	// If simulateBeforeRun is set to true, then transaction will be simulated
	// on the EVM before actual transaction will be send.
	Lift(ctx context.Context, addresses []common.Address, simulateBeforeRun bool) (*ethereum.Hash, error)
	// Drop sends transaction to the smart contract which invokes contract's
	// drop method, which removes given addresses from the feeders list (orcls).
	// If simulateBeforeRun is set to true, then transaction will be simulated
	// on the EVM before actual transaction will be send.
	Drop(ctx context.Context, addresses []common.Address, simulateBeforeRun bool) (*ethereum.Hash, error)
	// SetBar sends transaction to the smart contract which invokes contract's
	// setBar method, which sets bar variable (quorum).  If simulateBeforeRun is
	// set to true, then transaction will be simulated on the EVM before actual
	// transaction will be send.
	SetBar(ctx context.Context, bar *big.Int, simulateBeforeRun bool) (*ethereum.Hash, error)
}
