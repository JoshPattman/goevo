package goevo

// Agent is a container for a genotype and its fitness.
// The genotype can be of any type.
type Agent[T any] struct {
	Genotype T
	Fitness  float64
}

// NewAgent creates a new agent with the given genotype.
func NewAgent[T any](gt T) *Agent[T] {
	return &Agent[T]{Genotype: gt}
}
