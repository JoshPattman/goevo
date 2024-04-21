package neat

import (
	"math/rand"
	"slices"

	"github.com/JoshPattman/goevo"
	"golang.org/x/exp/maps"
)

type SimpleCrossoverStrategy struct{}

var _ goevo.CrossoverStrategy[*Genotype] = &SimpleCrossoverStrategy{}

type AsexualCrossoverStrategy struct{}

var _ goevo.CrossoverStrategy[*Genotype] = &AsexualCrossoverStrategy{}

// Crossover implements goevo.CrossoverStrategy.
func (s *SimpleCrossoverStrategy) Crossover(gs []*Genotype) *Genotype {
	if len(gs) != 2 {
		panic("expected 2 parents for simple crossover")
	}
	g, g2 := gs[0], gs[1]
	gc := &Genotype{
		g.maxSynapseValue,
		g.numInputs,
		g.numOutputs,
		slices.Clone(g.neuronOrder),
		maps.Clone(g.inverseNeuronOrder),
		maps.Clone(g.activations),
		maps.Clone(g.weights),
		maps.Clone(g.synapseEndpointLookup),
		maps.Clone(g.endpointSynapseLookup),
		slices.Clone(g.forwardSynapses),
		slices.Clone(g.backwardSynapses),
		slices.Clone(g.selfSynapses),
	}

	for sid, sw := range g2.weights {
		if _, ok := gc.weights[sid]; ok {
			if rand.Float64() > 0.5 {
				gc.weights[sid] = sw
			}
		}
	}

	return gc
}

// NumParents implements goevo.CrossoverStrategy.
func (s *SimpleCrossoverStrategy) NumParents() int {
	return 2
}

func (s *AsexualCrossoverStrategy) Crossover(gs []*Genotype) *Genotype {
	if len(gs) != 1 {
		panic("expected 1 parent for asexual crossover")
	}
	return goevo.Clone(gs[0])
}

func (s *AsexualCrossoverStrategy) NumParents() int {
	return 1
}
