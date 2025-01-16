package goevo

// Population is an interface for a population with genotypes with type T.
// It stores its genotypes wrapped in the [Agent] struct, to keep track of fitness.
// The population may also store a reference to a [Reproduction] and a [Selection]
// to be used in the [NextGeneration] method.
type Population[T any] interface {
	// NextGeneration returns the population resulting from agents selected using this population's selection strategy
	// reproducing using this population's reproduction strategy.
	NextGeneration() Population[T]

	// All returns every [Agent] in the population.
	// This may have no particular order.
	All() []*Agent[T]
}

func NextGeneration[T any, U Population[T]](pop U) U {
	return pop.NextGeneration().(U)
}
