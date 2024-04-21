package goevo

// SelectionStrategy is a strategy for selecting an [Agent] from a slice.
// It acts on agents of type T.
type SelectionStrategy[T any] interface {
	// SetAgents caches the [Agent]s which this selection will use until it is called again.
	// This is called once per generation. You may wish to perform slow operations here such as sorting by fitness.
	SetAgents(agents []*Agent[T])
	// Select returns an [Agent] selected from the cached pool set by [SelectionStrategy.SetAgents].
	Select() *Agent[T]
}

// SelectN selects n [Agent]s from the given selection strategy, returning them in a slice.
func SelectN[T any](selection SelectionStrategy[T], n int) []*Agent[T] {
	agents := make([]*Agent[T], n)
	for i := range agents {
		agents[i] = selection.Select()
	}
	return agents
}

// SelectNGenotypes selects n genotypes from the given selection strategy, returning them in a slice.
func SelectNGenotypes[T any](selection SelectionStrategy[T], n int) []T {
	gts := make([]T, n)
	for i := range gts {
		gts[i] = selection.Select().Genotype
	}
	return gts
}
