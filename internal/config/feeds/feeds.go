package feeds

import (
	"errors"
	"fmt"

	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
)

type Feeds []string

var ErrInvalidEthereumAddress = errors.New("invalid ethereum address")

func (f *Feeds) Addresses() ([]ethereum.Address, error) {
	var addrs []ethereum.Address
	for _, addr := range *f {
		if !ethereum.IsHexAddress(addr) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidEthereumAddress, addr)
		}
		addrs = append(addrs, ethereum.HexToAddress(addr))
	}
	return addrs, nil
}
