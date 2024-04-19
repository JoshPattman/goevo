package goevo

// Selection is a strategy for selecting agents from a population.
// It acts on agents of type T.
type Selection[T any] interface {
	// SetAgents sets the agents to select from for this generation.
	SetAgents(agents []*Agent[T])
	// Select returns an agent selected from the population.
	Select() *Agent[T]
}
