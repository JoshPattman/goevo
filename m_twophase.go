package goevo

// twoPhaseReproduction is a [Reproduction] that first performs a [Crossover]
// and then a [Mutation] on the resulting child.
type twoPhaseReproduction[T any] struct {
	crossover Crossover[T]
	mutate    Mutation[T]
}

// NewTwoPhaseReproduction creates a new [twoPhaseReproduction] with the given [Crossover] and [Mutation].
func NewTwoPhaseReproduction[T any](crossover Crossover[T], mutate Mutation[T]) Reproduction[T] {
	if crossover == nil {
		panic("cannot have nil crossover")
	}
	if mutate == nil {
		panic("cannot have nil mutate")
	}
	return &twoPhaseReproduction[T]{
		crossover: crossover,
		mutate:    mutate,
	}
}

// Reproduce implements the [Reproduction] interface.
func (r *twoPhaseReproduction[T]) Reproduce(parents []T) T {
	if len(parents) != r.crossover.NumParents() {
		panic("incorrect number of parents")
	}
	child := r.crossover.Crossover(parents)
	r.mutate.Mutate(child)
	return child
}

// NumParents implements the [Reproduction] interface.
func (r *twoPhaseReproduction[T]) NumParents() int {
	return r.crossover.NumParents()
}
