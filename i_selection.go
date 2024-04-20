package goevo

// SelectionStrategy is a strategy for selecting agents from a population.
// It acts on agents of type T.
type SelectionStrategy[T any] interface {
	// SetAgents sets the agents to select from for this generation.
	// This is run once per generation. You may wish to perform slow operations here such as sorting by fitness.
	SetAgents(agents []*Agent[T])
	// Select returns an agent selected from the agents set by SetAgents.
	Select() *Agent[T]
}

func SelectN[T any](selection SelectionStrategy[T], n int) []*Agent[T] {
	agents := make([]*Agent[T], n)
	for i := range agents {
		agents[i] = selection.Select()
	}
	return agents
}

func SelectNGenotypes[T any](selection SelectionStrategy[T], n int) []T {
	gts := make([]T, n)
	for i := range gts {
		gts[i] = selection.Select().Genotype
	}
	return gts
}
