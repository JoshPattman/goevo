package goevo

type Population[T any] interface {
	// NextGeneration returns the next generation of the population.
	// The population should use the selection and reproduction strategies it has stored to determine the next generation.
	NextGeneration() Population[T]

	// All returns all agents in the population.
	All() []*Agent[T]
}
