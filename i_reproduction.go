package goevo

// Mutation is an interface for a mutation strategy on a genotype with type T.
type Mutation[T any] interface {
	// Mutate performs a mutation in-place with this strategy on the given genotype
	Mutate(T)
}

// CrossoverStrategy is an interface for a crossover strategy on a genotype with type T.
type CrossoverStrategy[T any] interface {
	// Crossover performs a crossover with this strategy on the given genotypes.
	// It can combine any number of genotypes (for example 1 for asexual, 2 for sexual, n for averaging of multiple?)
	Crossover([]T) T
	// NumParents returns the number of parents required for this crossover strategy
	NumParents() int
}

// Reproduction is an interface for a reproduction strategy on a genotype with type T.
// Most of the time, this will be a [TwoPhaseReproduction], however it
// is possible to imlement a custom one for more complex behaviour.
type Reproduction[T any] interface {
	// Reproduce takes a set of parent genotypes and returns a child genotype.
	Reproduce([]T) T
	// NumParents returns the number of parents required for this reproduction strategy
	NumParents() int
}
