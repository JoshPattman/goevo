package goevo

import (
	"sort"

	"math/rand/v2"
)

// ArrayGenotype is a genotype that is a slice of values.
type ArrayGenotype[T any] struct {
	values []T
}

func NewArrayGenotype[T any](length int, gen Generator[T]) *ArrayGenotype[T] {
	if length <= 0 {
		panic("must have length > 1")
	}
	if gen == nil {
		panic("must have non-nil generator")
	}
	vals := make([]T, length)
	for i := range vals {
		vals[i] = gen.Next()
	}
	return &ArrayGenotype[T]{
		values: vals,
	}
}

func (g *ArrayGenotype[T]) Len() int {
	return len(g.values)
}

func (g *ArrayGenotype[T]) At(i int) T {
	return g.values[i]
}

func (g *ArrayGenotype[T]) Set(i int, v T) {
	g.values[i] = v
}

// Clone returns a new genotype that is a copy of this genotype.
func (g *ArrayGenotype[T]) Clone() any {
	clone := make([]T, len(g.values))
	copy(clone, g.values)
	return &ArrayGenotype[T]{values: clone}
}

// arrayCrossoverUniform is a crossover strategy that selects each gene from one of the parents with equal probability.
// The location of a gene has no effect on the probability of it being selected from either parent.
// It requires two parents.
type arrayCrossoverUniform[T any] struct{}

func NewArrayCrossoverUniform[T any]() Crossover[*ArrayGenotype[T]] {
	return &arrayCrossoverUniform[T]{}
}

// Crossover implements CrossoverStrategy.
func (p *arrayCrossoverUniform[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
	if len(gs) != 2 {
		panic("PointCrossoverStrategy requires exactly 2 parents")
	}
	pa, pb := gs[0], gs[1]
	if len(pa.values) != len(pb.values) {
		panic("genotypes must have the same length for PointCrossoverStrategy")
	}
	child := make([]T, len(pa.values))
	for i := range child {
		if rand.Float64() < 0.5 {
			child[i] = pa.values[i]
		} else {
			child[i] = pb.values[i]
		}
	}
	return &ArrayGenotype[T]{values: child}
}

// NumParents implements CrossoverStrategy.
func (p *arrayCrossoverUniform[T]) NumParents() int {
	return 2
}

// arrayCrossoverAsexual is a crossover strategy that clones the parent.
// It only requires one parent.
type arrayCrossoverAsexual[T any] struct{}

func NewArrayCrossoverAsexual[T any]() Crossover[*ArrayGenotype[T]] {
	return &arrayCrossoverAsexual[T]{}
}

// Crossover implements CrossoverStrategy.
func (p *arrayCrossoverAsexual[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
	if len(gs) != 1 {
		panic("AsexualCrossoverStrategy requires exactly 1 parent")
	}
	return Clone(gs[0])
}

// NumParents implements CrossoverStrategy.
func (p *arrayCrossoverAsexual[T]) NumParents() int {
	return 1
}

// arrayCrossoverKPoint is a crossover strategy that selects K locations in the genome to switch parents.
// It requires two parents.
type arrayCrossoverKPoint[T any] struct {
	k int
}

func NewArrayCrossoverKPoint[T any](k int) Crossover[*ArrayGenotype[T]] {
	if k < 0 {
		panic("k must be > 0")
	}
	return &arrayCrossoverKPoint[T]{
		k: k,
	}
}

func (p *arrayCrossoverKPoint[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
	if len(gs) != 2 {
		panic("KPointCrossoverStrategy requires exactly 2 parents")
	}
	pa, pb := gs[0], gs[1]
	if len(pa.values) != len(pb.values) {
		panic("genotypes must have the same length for KPointCrossoverStrategy")
	}
	crossoverPoints := make([]int, p.k)
	for i := 0; i < p.k; i++ {
		crossoverPoints[i] = rand.N(len(pa.values))
	}
	sort.Ints(crossoverPoints)
	child := make([]T, len(pa.values))
	fromParentA := rand.Float64() < 0.5
	currentCrossoverPoint := 0
	for i := range child {
		if currentCrossoverPoint < len(crossoverPoints) && crossoverPoints[currentCrossoverPoint] == i {
			fromParentA = !fromParentA
			currentCrossoverPoint++
		}
		if fromParentA {
			child[i] = pa.values[i]
		} else {
			child[i] = pb.values[i]
		}
	}
	return &ArrayGenotype[T]{values: child}
}

func (p *arrayCrossoverKPoint[T]) NumParents() int {
	return 2
}

type arrayMutationGenerator[T any] struct {
	combine func(old, new T) T
	gen     Generator[T]
	chance  float64
}

func NewArrayMutationGeneratorAdd[T numberType](gen Generator[T], chance float64) Mutation[*ArrayGenotype[T]] {
	combine := func(old, new T) T {
		return old + new
	}
	return NewArrayMutationGenerator(gen, combine, chance)
}

func NewArrayMutationGeneratorReplace[T any](gen Generator[T], chance float64) Mutation[*ArrayGenotype[T]] {
	combine := func(_, new T) T {
		return new
	}
	return NewArrayMutationGenerator(gen, combine, chance)
}

func NewArrayMutationGenerator[T any](gen Generator[T], combine func(old, new T) T, chance float64) Mutation[*ArrayGenotype[T]] {
	if gen == nil {
		panic("cannot have nil generator")
	}
	if combine == nil {
		panic("cannot have nil combine")
	}
	if chance < 0 {
		panic("cannot have chance < 0")
	}
	return &arrayMutationGenerator[T]{
		combine: combine,
		gen:     gen,
		chance:  chance,
	}
}

// Mutate implements Mutation.
func (m *arrayMutationGenerator[T]) Mutate(g *ArrayGenotype[T]) {
	for i := range g.Len() {
		g.values[i] = m.combine(g.values[i], m.gen.Next())
	}
}
