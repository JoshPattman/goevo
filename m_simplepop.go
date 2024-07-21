package goevo

var _ Population[int] = &SimplePopulation[int]{}

// SimplePopulation has a single species, and generates the entire next generation by selcting and breeding from the previous one.
type SimplePopulation[T any] struct {
	Agents       []*Agent[T]
	Selection    Selection[T]
	Reproduction Reproduction[T]
}

// NewSimplePopulation creates a new SimplePopulation with n agents, each with a new genotype created by newGenotype.
func NewSimplePopulation[T any](newGenotype func() T, n int, selection Selection[T], reproduction Reproduction[T]) *SimplePopulation[T] {
	agents := make([]*Agent[T], n)
	for i := range agents {
		agents[i] = NewAgent(newGenotype())
	}
	return &SimplePopulation[T]{
		Agents:       agents,
		Selection:    selection,
		Reproduction: reproduction,
	}
}

// NextGeneration creates a new SimplePopulation from the current one, using the given selection and reproduction strategies.
func (p *SimplePopulation[T]) NextGeneration() Population[T] {
	p.Selection.SetAgents(p.Agents)
	return NewSimplePopulation(func() T {
		parents := SelectNGenotypes(p.Selection, p.Reproduction.NumParents())
		return p.Reproduction.Reproduce(parents)
	}, len(p.Agents), p.Selection, p.Reproduction)
}

// Agents returns the agents in the population.
//
// TODO(change this to an iterator once they get added to the language, as this will increase performance in other cases)
func (p *SimplePopulation[T]) All() []*Agent[T] {
	return p.Agents
}
