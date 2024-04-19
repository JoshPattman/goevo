package goevo

// Selection is a strategy for selecting agents from a population.
// It acts on agents of type T.
type Selection[T any] interface {
	// SetAgents sets the agents to select from for this generation.
	SetAgents(agents []*Agent[T])
	// Select returns an agent selected from the population.
	Select() *Agent[T]
}

// Reproduction is an interface for the reproduction of two parents to create a child
type Reproduction[T any] interface {
	// Reproduce creates a new genotype from the two parents, where the first parent is fitter
	Reproduce(a, b T) T
}

// Forwarder is an interface for somthing that can take a set of inputs ([]float64) and return a set of outputs.
type Forwarder interface {
	Forward([]float64) []float64
}

type GeneticDistance[T any] interface {
	DistanceBetween(a, b T) float64
}
