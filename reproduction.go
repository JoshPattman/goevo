package goevo

// Implementations
var _ ReproductionStrategy[int] = &CrossoverMutateReproduction[int]{}

// CrossoverMutateReproduction is a [ReproductionStrategy] that first performs a [CrossoverStrategy]
// and then a [MutationStrategy] on the resulting child.
type CrossoverMutateReproduction[T any] struct {
	Crossover CrossoverStrategy[T]
	Mutate    MutationStrategy[T]
}

// NewCrossoverMutateReproduction creates a new [CrossoverMutateReproduction] with the given [CrossoverStrategy] and [MutationStrategy].
func NewCrossoverMutateReproduction[T any](crossover CrossoverStrategy[T], mutate MutationStrategy[T]) *CrossoverMutateReproduction[T] {
	return &CrossoverMutateReproduction[T]{
		Crossover: crossover,
		Mutate:    mutate,
	}
}

// Reproduce implements the [ReproductionStrategy] interface.
func (r *CrossoverMutateReproduction[T]) Reproduce(parents []T) T {
	if len(parents) != r.Crossover.NumParents() {
		panic("incorrect number of parents")
	}
	child := r.Crossover.Crossover(parents)
	r.Mutate.Mutate(child)
	return child
}

// NumParents implements the [ReproductionStrategy] interface.
func (r *CrossoverMutateReproduction[T]) NumParents() int {
	return r.Crossover.NumParents()
}
