package goevo

// Reproduction is an interface for the reproduction of two parents to create a child
type Reproduction[T any] interface {
	// Reproduce creates a new genotype from the two parents, where the first parent is fitter
	Reproduce(a, b T) T
}
