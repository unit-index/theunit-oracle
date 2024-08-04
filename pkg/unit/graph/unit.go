package graph

import (
	"github.com/toknowwhy/theunit-oracle/pkg/unit/graph/feeder"
)

type Unit struct {
	feeder *feeder.Feeder
}

func NewUnit(f *feeder.Feeder) *Unit {
	return &Unit{feeder: f}
}

func (u *Unit) CSupply() {

}
