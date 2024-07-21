package goevo

import (
	"sort"

	"math/rand/v2"
)

var _ Cloneable = ArrayGenotype[int]{}

// ArrayGenotype is a genotype that is a slice of values.
type ArrayGenotype[T any] struct {
	Values []T
}

// NewFloatArrayGenotype creates a new genotype with the given size,
// where each value is a random float (of type T) with a normal distribution with the given standard deviation.
func NewFloatArrayGenotype[T floatType](size int, std T) *ArrayGenotype[T] {
	gt := make([]T, size)
	for i := range gt {
		gt[i] = T(rand.NormFloat64()) * std
	}
	return &ArrayGenotype[T]{Values: gt}
}

// NewRuneArrayGenotype creates a new genotype with the given size,
// where each value is a random rune from the given runeset.
func NewRuneArrayGenotype(size int, runeset []rune) *ArrayGenotype[rune] {
	gt := make([]rune, size)
	for i := range gt {
		gt[i] = runeset[rand.N(len(runeset))]
	}
	return &ArrayGenotype[rune]{Values: gt}
}

// NewBoolArrayGenotype creates a new genotype with the given size,
// where each value is a random boolean.
func NewBoolArrayGenotype(size int) *ArrayGenotype[bool] {
	gt := make([]bool, size)
	for i := range gt {
		gt[i] = rand.Float64() < 0.5
	}
	return &ArrayGenotype[bool]{Values: gt}
}

// Clone returns a new genotype that is a copy of this genotype.
func (g ArrayGenotype[T]) Clone() any {
	clone := make([]T, len(g.Values))
	copy(clone, g.Values)
	return &ArrayGenotype[T]{Values: clone}
}

var _ CrossoverStrategy[*ArrayGenotype[int]] = &ArrayCrossoverUniform[int]{}

// ArrayCrossoverUniform is a crossover strategy that selects each gene from one of the parents with equal probability.
// The location of a gene has no effect on the probability of it being selected from either parent.
// It requires two parents.
type ArrayCrossoverUniform[T any] struct{}

// Crossover implements CrossoverStrategy.
func (p *ArrayCrossoverUniform[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
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
	return &ArrayGenotype[T]{Values: child}
}

// NumParents implements CrossoverStrategy.
func (p *ArrayCrossoverUniform[T]) NumParents() int {
	return 2
}

var _ CrossoverStrategy[*ArrayGenotype[bool]] = &ArrayCrossoverAsexual[bool]{}

// ArrayCrossoverAsexual is a crossover strategy that clones the parent.
// It only requires one parent.
type ArrayCrossoverAsexual[T any] struct{}

// Crossover implements CrossoverStrategy.
func (p *ArrayCrossoverAsexual[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
	if len(gs) != 1 {
		panic("AsexualCrossoverStrategy requires exactly 1 parent")
	}
	return Clone(gs[0])
}

// NumParents implements CrossoverStrategy.
func (p *ArrayCrossoverAsexual[T]) NumParents() int {
	return 1
}

var _ CrossoverStrategy[*ArrayGenotype[int]] = &ArrayCrossoverKPoint[int]{}

// ArrayCrossoverKPoint is a crossover strategy that selects K locations in the genome to switch parents.
// It requires two parents.
type ArrayCrossoverKPoint[T any] struct {
	K int
}

func (p *ArrayCrossoverKPoint[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
	if len(gs) != 2 {
		panic("KPointCrossoverStrategy requires exactly 2 parents")
	}
	pa, pb := gs[0], gs[1]
	if len(pa.Values) != len(pb.Values) {
		panic("genotypes must have the same length for KPointCrossoverStrategy")
	}
	crossoverPoints := make([]int, p.K)
	for i := 0; i < p.K; i++ {
		crossoverPoints[i] = rand.N(len(pa.Values))
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
	return &ArrayGenotype[T]{Values: child}
}

func (p *ArrayCrossoverKPoint[T]) NumParents() int {
	return 2
}

var _ Mutation[*ArrayGenotype[bool]] = &ArrayMutationRandomBool{}

type ArrayMutationRandomBool struct {
	// The probability of mutating each locus
	MutateProbability float64
}

func (s *ArrayMutationRandomBool) Mutate(gt *ArrayGenotype[bool]) {
	for i := range gt.Values {
		if rand.Float64() < s.MutateProbability {
			gt.Values[i] = !gt.Values[i]
		}
	}
}

var _ Mutation[*ArrayGenotype[float64]] = &ArrayMutationStd[float64]{}

type ArrayMutationStd[T floatType] struct {
	// The probability of mutating each locus
	MutateProbability T
	// The standard deviation for the mutation
	MutateStd T
}

func (s *ArrayMutationStd[T]) Mutate(gt *ArrayGenotype[T]) {
	for i := range gt.Values {
		if T(rand.Float64()) < s.MutateProbability {
			gt.Values[i] += T(rand.NormFloat64()) * s.MutateStd
		}
	}
}

var _ Mutation[*ArrayGenotype[rune]] = &ArrayMutationRandomRune{}

type ArrayMutationRandomRune struct {
	// The probability of mutating each locus
	MutateProbability float64
	// The standard deviation for the mutation
	Runeset []rune
}

func (s *ArrayMutationRandomRune) Mutate(gt *ArrayGenotype[rune]) {
	for i := range gt.Values {
		if rand.Float64() < s.MutateProbability {
			gt.Values[i] = s.Runeset[rand.N(len(s.Runeset))]
		}
	}
}
