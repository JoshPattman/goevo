package goevo

// Reproduction is an interface for the reproduction of n parents to create a child
type Reproduction[T any] interface {
	// Reproduce creates a new genotype from the n parents. The parents are NOT ordered by fitness.
	Reproduce(agents []T) T
	// NumParents returns the number of parents required for reproduction
	NumParents() int
}
