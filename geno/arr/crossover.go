package arr

import (
	"math/rand"
	"sort"

	"github.com/JoshPattman/goevo"
)

var _ goevo.CrossoverStrategy[*Genotype[int]] = &UniformCrossoverStrategy[int]{}
var _ goevo.CrossoverStrategy[*Genotype[bool]] = &AsexualCrossoverStrategy[bool]{}
var _ goevo.CrossoverStrategy[*Genotype[int]] = &KPointCrossoverStrategy[int]{}

// UniformCrossoverStrategy is a crossover strategy that selects each gene from one of the parents with equal probability.
// The location of a gene has no effect on the probability of it being selected from either parent.
// It requires two parents.
type UniformCrossoverStrategy[T any] struct{}

// AsexualCrossoverStrategy is a crossover strategy that clones the parent.
// It only requires one parent.
type AsexualCrossoverStrategy[T any] struct{}

// KPointCrossoverStrategy is a crossover strategy that selects K locations in the genome to switch parents.
// It requires two parents.
type KPointCrossoverStrategy[T any] struct {
	K int
}

// Crossover implements goevo.CrossoverStrategy.
func (p *UniformCrossoverStrategy[T]) Crossover(gs []*Genotype[T]) *Genotype[T] {
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
func (p *UniformCrossoverStrategy[T]) NumParents() int {
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

func (p *KPointCrossoverStrategy[T]) Crossover(gs []*Genotype[T]) *Genotype[T] {
	if len(gs) != 2 {
		panic("KPointCrossoverStrategy requires exactly 2 parents")
	}
	pa, pb := gs[0], gs[1]
	if len(pa.Values) != len(pb.Values) {
		panic("genotypes must have the same length for KPointCrossoverStrategy")
	}
	crossoverPoints := make([]int, p.K)
	for i := 0; i < p.K; i++ {
		crossoverPoints[i] = rand.Intn(len(pa.Values))
	}
	sort.Ints(crossoverPoints)
	child := make([]T, len(pa.Values))
	fromParentA := rand.Float64() < 0.5
	currentCrossoverPoint := 0
	for i := range child {
		if currentCrossoverPoint < len(crossoverPoints) && crossoverPoints[currentCrossoverPoint] == i {
			fromParentA = !fromParentA
			currentCrossoverPoint++
		}
		if fromParentA {
			child[i] = pa.Values[i]
		} else {
			child[i] = pb.Values[i]
		}
	}
	return &Genotype[T]{Values: child}
}

func (p *KPointCrossoverStrategy[T]) NumParents() int {
	return 2
}
