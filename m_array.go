package goevo

import (
	"fmt"
	"sort"

	"math/rand/v2"
)

// ArrayGenotype is a genotype that is a slice of values.
type ArrayGenotype[T any] struct {
	Values []T
}

// Validate implements Validateable.
func (g *ArrayGenotype[T]) Validate() error {
	if g.Values == nil {
		return fmt.Errorf("an array genotype may not have nil values array")
	}
	return nil
}

// Clone implements Cloneable.
func (g *ArrayGenotype[T]) Clone() any {
	MustValidate(g)
	clone := make([]T, len(g.Values))
	copy(clone, g.Values)
	return &ArrayGenotype[T]{Values: clone}
}

// ArrayFactoryGenerator creates [ArrayGenotype]s
// using the provided [Generator] with the given length.
type ArrayFactoryGenerator[T any] struct {
	Length    int
	Generator Generator[T]
}

// New implements ValidateableFactory.
func (a *ArrayFactoryGenerator[T]) New() *ArrayGenotype[T] {
	MustValidate(a)
	gt := make([]T, a.Length)
	for i := range gt {
		gt[i] = a.Generator.Next()
	}
	return &ArrayGenotype[T]{Values: gt}
}

// Validate implements ValidateableFactory.
func (a *ArrayFactoryGenerator[T]) Validate() error {
	if a.Length < 0 {
		return fmt.Errorf("cannot create array genotypes with negative (%v) size", a.Length)
	}
	if a.Generator == nil {
		return fmt.Errorf("cannot use an ArrayFactoryGenerator with a nil generator")
	}
	return nil
}

// ArrayCrossoverUniform is a crossover strategy that selects each gene from one of the parents with equal probability.
// The location of a gene has no effect on the probability of it being selected from either parent.
// It requires two parents.
type ArrayCrossoverUniform[T any] struct{}

// Crossover implements CrossoverStrategy.
func (p *ArrayCrossoverUniform[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
	MustValidateAll(gs...)
	MustValidate(p)
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

func (*ArrayCrossoverUniform[T]) Validate() error { return nil }

// ArrayCrossoverAsexual is a crossover strategy that clones the parent.
// It only requires one parent.
type ArrayCrossoverAsexual[T any] struct{}

// Crossover implements CrossoverStrategy.
func (p *ArrayCrossoverAsexual[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
	MustValidateAll(gs...)
	MustValidate(p)
	if len(gs) != 1 {
		panic("AsexualCrossoverStrategy requires exactly 1 parent")
	}
	return Clone(gs[0])
}

// NumParents implements CrossoverStrategy.
func (p *ArrayCrossoverAsexual[T]) NumParents() int {
	return 1
}

func (*ArrayCrossoverAsexual[T]) Validate() error { return nil }

// ArrayCrossoverKPoint is a crossover strategy that selects K locations in the genome to switch parents.
// It requires two parents.
type ArrayCrossoverKPoint[T any] struct {
	K int
}

func (p *ArrayCrossoverKPoint[T]) Crossover(gs []*ArrayGenotype[T]) *ArrayGenotype[T] {
	MustValidateAll(gs...)
	MustValidate(p)
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

func (p *ArrayCrossoverKPoint[T]) Validate() error {
	if p.K < 0 {
		return fmt.Errorf("cannot set K < 0")
	}
	return nil
}

type ArrayMutationRandomBool struct {
	// The probability of mutating each locus
	MutateProbability float64
}

func (s *ArrayMutationRandomBool) Mutate(gt *ArrayGenotype[bool]) {
	MustValidate(gt)
	MustValidate(s)
	for i := range gt.Values {
		if rand.Float64() < s.MutateProbability {
			gt.Values[i] = !gt.Values[i]
		}
	}
}

func (p *ArrayMutationRandomBool) Validate() error {
	if p.MutateProbability < 0 || p.MutateProbability > 1 {
		return fmt.Errorf("cannot set mutation probability out of range 0-1")
	}
	return nil
}

type ArrayMutationStd[T floatType] struct {
	// The probability of mutating each locus
	MutateProbability T
	// The standard deviation for the mutation
	MutateStd T
}

func (s *ArrayMutationStd[T]) Mutate(gt *ArrayGenotype[T]) {
	MustValidate(gt)
	MustValidate(s)
	for i := range gt.Values {
		if T(rand.Float64()) < s.MutateProbability {
			gt.Values[i] += T(rand.NormFloat64()) * s.MutateStd
		}
	}
}

func (p *ArrayMutationStd[T]) Validate() error {
	if p.MutateProbability < 0 || p.MutateProbability > 1 {
		return fmt.Errorf("cannot set mutation probability out of range 0-1")
	}
	if p.MutateStd < 0 {
		return fmt.Errorf("cannot set mutate std < 0")
	}
	return nil
}

type ArrayMutationRandomRune struct {
	// The probability of mutating each locus
	MutateProbability float64
	// The standard deviation for the mutation
	Runeset []rune
}

func (s *ArrayMutationRandomRune) Mutate(gt *ArrayGenotype[rune]) {
	MustValidate(gt)
	MustValidate(s)
	for i := range gt.Values {
		if rand.Float64() < s.MutateProbability {
			gt.Values[i] = s.Runeset[rand.N(len(s.Runeset))]
		}
	}
}

func (p *ArrayMutationRandomRune) Validate() error {
	if p.MutateProbability < 0 || p.MutateProbability > 1 {
		return fmt.Errorf("cannot set mutation probability out of range 0-1")
	}
	if len(p.Runeset) == 0 {
		return fmt.Errorf("cannot have a runeset of length 0")
	}
	return nil
}
