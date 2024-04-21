package floatarr

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

var _ goevo.CrossoverStrategy[*Genotype] = &PointCrossoverStrategy{}
var _ goevo.CrossoverStrategy[*Genotype] = &AsexualCrossoverStrategy{}

type PointCrossoverStrategy struct{}
type AsexualCrossoverStrategy struct{}

// Crossover implements goevo.CrossoverStrategy.
func (p *PointCrossoverStrategy) Crossover(gs []*Genotype) *Genotype {
	if len(gs) != 2 {
		panic("PointCrossoverStrategy requires exactly 2 parents")
	}
	pa, pb := gs[0], gs[1]
	if len(pa.Values) != len(pb.Values) {
		panic("genotypes must have the same length for PointCrossoverStrategy")
	}
	child := make([]float64, len(pa.Values))
	for i := range child {
		if rand.Float64() < 0.5 {
			child[i] = pa.Values[i]
		} else {
			child[i] = pb.Values[i]
		}
	}
	return &Genotype{Values: child}
}

// NumParents implements goevo.CrossoverStrategy.
func (p *PointCrossoverStrategy) NumParents() int {
	return 2
}

// Crossover implements goevo.CrossoverStrategy.
func (p *AsexualCrossoverStrategy) Crossover(gs []*Genotype) *Genotype {
	if len(gs) != 1 {
		panic("AsexualCrossoverStrategy requires exactly 1 parent")
	}
	return goevo.Clone(gs[0])
}

// NumParents implements goevo.CrossoverStrategy.
func (p *AsexualCrossoverStrategy) NumParents() int {
	return 1
}
