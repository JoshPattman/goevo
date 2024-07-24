package goevo

// Implementations
var _ Reproduction[int] = &TwoPhaseReproduction[int]{}

// TwoPhaseReproduction is a [Reproduction] that first performs a [Crossover]
// and then a [Mutation] on the resulting child.
type TwoPhaseReproduction[T any] struct {
	Crossover Crossover[T]
	Mutate    Mutation[T]
}

// NewTwoPhaseReproduction creates a new [TwoPhaseReproduction] with the given [Crossover] and [Mutation].
func NewTwoPhaseReproduction[T any](crossover Crossover[T], mutate Mutation[T]) *TwoPhaseReproduction[T] {
	return &TwoPhaseReproduction[T]{
		Crossover: crossover,
		Mutate:    mutate,
	}
}

// Reproduce implements the [Reproduction] interface.
func (r *TwoPhaseReproduction[T]) Reproduce(parents []T) T {
	if len(parents) != r.Crossover.NumParents() {
		panic("incorrect number of parents")
	}
	child := r.Crossover.Crossover(parents)
	r.Mutate.Mutate(child)
	return child
}

// NumParents implements the [Reproduction] interface.
func (r *TwoPhaseReproduction[T]) NumParents() int {
	return r.Crossover.NumParents()
}
