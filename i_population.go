package goevo

// Population is an interface for a population with genotypes with type T.
// It stores its genotypes wrapped in the [Agent] struct, to keep track of fitness.
// The population may also store a reference to a [Reproduction] and a [Selection]
// to be used in the [NextGeneration] method.
// Type T is usually a pointer type (for example, T=*NEATGenotype)
type Population[T any] interface {
	// NextGeneration returns the population resulting from agents selected using this population's selection strategy
	// reproducing using this population's reproduction strategy.
	// It takes in a function (that can be nil) that is run on every
	// T that was NOT reused (memory level) in the new generation.
	// This means you can use this function as a place to return old agents to a memory pool, etc.
	// If you pass nil here, the go GC will just have to take care of the objects.
	NextGeneration(recycle func(T)) Population[T]

	// All returns every [Agent] in the population.
	// This may have no particular order.
	All() []*Agent[T]
}
