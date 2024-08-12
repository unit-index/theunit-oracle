//  Copyright (C) 2020 Maker Ecosystem Growth Holdings, INC.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as
//  published by the Free Software Foundation, either version 3 of the
//  License, or (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package nodes

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/toknowwhy/theunit-oracle/pkg/unit"
)

type ErrIncompatibleToken struct {
	Given    unit.Token
	Expected unit.Token
}

func (e ErrIncompatibleToken) Error() string {
	return fmt.Sprintf(
		"a price with different token ignested to the OriginNode, %s given but %s was expected",
		e.Given,
		e.Expected,
	)
}

type IncompatibleOriginErr struct {
	Given    string
	Expected string
}

func (e IncompatibleOriginErr) Error() string {
	return fmt.Sprintf(
		"a price from different origin ignested to the OriginNode, %s given but %s was expected",
		e.Given,
		e.Expected,
	)
}

type ErrCSupplyTTLExpired struct {
	CSupply OriginCSupply
	TTL     time.Duration
}

func (e ErrCSupplyTTLExpired) Error() string {
	return fmt.Sprintf(
		"the price TTL for the token %s expired",
		e.CSupply.Token,
	)
}

// OriginNode contains a Price fetched directly from an origin.
type OriginNode struct {
	mu sync.RWMutex

	originToken OriginToken
	cSupply     OriginCSupply
	minTTL      time.Duration
	maxTTL      time.Duration
}

func NewOriginNode(originToken OriginToken, minTTL time.Duration, maxTTL time.Duration) *OriginNode {
	return &OriginNode{
		originToken: originToken,
		minTTL:      minTTL,
		maxTTL:      maxTTL,
	}
}

// OriginPair implements the Feedable interface.
func (n *OriginNode) OriginToken() OriginToken {
	return n.originToken
}

// Ingest implements Feedable interface.
func (n *OriginNode) Ingest(cSupply OriginCSupply) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	var err error
	if !cSupply.Token.Equal(n.originToken.Token) {
		err = multierror.Append(err, ErrIncompatibleToken{
			Given:    cSupply.Token,
			Expected: n.originToken.Token,
		})
	}

	if cSupply.Origin != n.originToken.Origin {
		err = multierror.Append(err, IncompatibleOriginErr{
			Given:    cSupply.Origin,
			Expected: n.originToken.Origin,
		})
	}

	if err == nil {
		n.cSupply = cSupply
	}

	return err
}

// MinTTL implements the Feedable interface.
func (n *OriginNode) MinTTL() time.Duration {
	return n.minTTL
}

// MaxTTL implements the Feedable interface.
func (n *OriginNode) MaxTTL() time.Duration {
	return n.maxTTL
}

// Expired implements the Feedable interface.
func (n *OriginNode) Expired() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.expired()
}

// Price implements the Feedable interface.
func (n *OriginNode) CSupply() OriginCSupply {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.cSupply.Error == nil {
		if n.expired() {
			n.cSupply.Error = ErrCSupplyTTLExpired{
				CSupply: n.cSupply,
				TTL:     n.maxTTL,
			}
		}
	}

	return n.cSupply
}

// Children implements the Node interface.
func (n *OriginNode) Children() []Node {
	return []Node{}
}

func (n *OriginNode) expired() bool {
	return n.cSupply.Time.Before(time.Now().Add(-1 * n.MaxTTL()))
}
