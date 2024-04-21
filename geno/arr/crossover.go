package arr

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

var _ goevo.CrossoverStrategy[*Genotype[int]] = &PointCrossoverStrategy[int]{}
var _ goevo.CrossoverStrategy[*Genotype[bool]] = &AsexualCrossoverStrategy[bool]{}

type PointCrossoverStrategy[T any] struct{}
type AsexualCrossoverStrategy[T any] struct{}

// Crossover implements goevo.CrossoverStrategy.
func (p *PointCrossoverStrategy[T]) Crossover(gs []*Genotype[T]) *Genotype[T] {
	if len(gs) != 2 {
		panic("PointCrossoverStrategy requires exactly 2 parents")
	}
	pa, pb := gs[0], gs[1]
	if len(pa.Values) != len(pb.Values) {
		panic("genotypes must have the same length for PointCrossoverStrategy")
	}
	child := make([]T, len(pa.Values))
	for i := range child {
		if rand.Float64() < 0.5 {
			child[i] = pa.Values[i]
		} else {
			child[i] = pb.Values[i]
		}
	}
	return &Genotype[T]{Values: child}
}

// NumParents implements goevo.CrossoverStrategy.
func (p *PointCrossoverStrategy[T]) NumParents() int {
	return 2
}

// Crossover implements goevo.CrossoverStrategy.
func (p *AsexualCrossoverStrategy[T]) Crossover(gs []*Genotype[T]) *Genotype[T] {
	if len(gs) != 1 {
		panic("AsexualCrossoverStrategy requires exactly 1 parent")
	}
	return goevo.Clone(gs[0])
}

// NumParents implements goevo.CrossoverStrategy.
func (p *AsexualCrossoverStrategy[T]) NumParents() int {
	return 1
}
