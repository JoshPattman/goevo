package goevo

var _ Reproduction[int] = &CrossoverMutateReproduction[int]{}

type CrossoverMutateReproduction[T any] struct {
	Crossover CrossoverStrategy[T]
	Mutate    MutationStrategy[T]
}

func NewCrossoverMutateReproduction[T any](crossover CrossoverStrategy[T], mutate MutationStrategy[T]) *CrossoverMutateReproduction[T] {
	return &CrossoverMutateReproduction[T]{
		Crossover: crossover,
		Mutate:    mutate,
	}
}

func (r *CrossoverMutateReproduction[T]) Reproduce(parents []T) T {
	if len(parents) != r.Crossover.NumParents() {
		panic("incorrect number of parents")
	}
	child := r.Crossover.Crossover(parents)
	r.Mutate.Mutate(child)
	return child
}

func (r *CrossoverMutateReproduction[T]) NumParents() int {
	return r.Crossover.NumParents()
}
