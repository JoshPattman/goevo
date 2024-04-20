package goevo

// Selection is a strategy for selecting agents from a population.
// It acts on agents of type T.
type Selection[T any] interface {
	// SetAgents sets the agents to select from for this generation.
	// This is run once per generation. You may wish to perform slow operations here such as sorting by fitness.
	SetAgents(agents []*Agent[T])
	// Select returns an agent selected from the agents set by SetAgents.
	Select() *Agent[T]
}
