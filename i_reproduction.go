package goevo

type MutationStrategy[T any] interface {
	// Mutate performs a mutation with this strategy on the given genotype
	Mutate(T)
}

type CrossoverStrategy[T any] interface {
	// Crossover performs a crossover with this strategy on the given genotypes.
	// It can combine any number of genotypes (for example 1 for asexual, 2 for sexual, n for averaging of multiple?)
	Crossover([]T) T
	// NumParents returns the number of parents required for this crossover strategy
	NumParents() int
}

type Reproduction[T any] interface {
	Reproduce([]T) T
	NumParents() int
}
