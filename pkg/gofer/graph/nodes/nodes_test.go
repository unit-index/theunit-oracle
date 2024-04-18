package nodes

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/toknowwhy/theunit-oracle/pkg/gofer"
)

func TestDetectCycles(t *testing.T) {
	p := gofer.Pair{Base: "A", Quote: "B"}

	// Non cyclic graph:
	nonCyclic := NewMedianAggregatorNode(p, 0)
	nonCyclicC1 := NewOriginNode(OriginPair{Origin: "a", Pair: p}, 0, 0)
	nonCyclicC2 := NewOriginNode(OriginPair{Origin: "b", Pair: p}, 0, 0)
	nonCyclicC3 := NewMedianAggregatorNode(p, 0)
	nonCyclic.AddChild(nonCyclicC1)
	nonCyclic.AddChild(nonCyclicC2)
	nonCyclic.AddChild(nonCyclicC3)
	nonCyclicC3.AddChild(nonCyclicC1)
	nonCyclicC3.AddChild(nonCyclicC2)

	// Cyclic graph:
	cyclic := NewMedianAggregatorNode(p, 0)
	cyclicC1 := NewOriginNode(OriginPair{Origin: "a", Pair: p}, 0, 0)
	cyclicC2 := NewOriginNode(OriginPair{Origin: "b", Pair: p}, 0, 0)
	cyclicC3 := NewMedianAggregatorNode(p, 0)
	cyclic.AddChild(cyclicC1)
	cyclic.AddChild(cyclicC2)
	cyclic.AddChild(cyclicC3)
	cyclicC3.AddChild(cyclicC1)
	cyclicC3.AddChild(cyclic)

	// Graph with references to the same aggregator nodes:
	r := NewMedianAggregatorNode(p, 0)
	c1 := NewMedianAggregatorNode(p, 0)
	c2 := NewMedianAggregatorNode(p, 0)
	r.AddChild(c1)
	r.AddChild(c2)
	r.AddChild(nonCyclic)
	c1.AddChild(nonCyclic)
	c2.AddChild(nonCyclic)
	c2.AddChild(cyclic)

	assert.Len(t, DetectCycle(nonCyclic), 0)
	assert.Equal(t, []Node{cyclic, cyclicC3}, DetectCycle(cyclic))
	assert.Equal(t, []Node{r, c2, cyclic, cyclicC3}, DetectCycle(r))
}
