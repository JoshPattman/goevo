package goevo

// SimplePopulation has a single species, and generates the entire next generation by selcting and breeding from the previous one.
type SimplePopulation[T any] struct {
	agents []*Agent[T]
}

// NewSimplePopulation creates a new SimplePopulation with n agents, each with a new genotype created by newGenotype.
func NewSimplePopulation[T any](newGenotype func() T, n int) *SimplePopulation[T] {
	agents := make([]*Agent[T], n)
	for i := range agents {
		agents[i] = NewAgent(newGenotype())
	}
	return &SimplePopulation[T]{
		agents: agents,
	}
}

// NextGeneration creates a new SimplePopulation from the current one, using the given selection and reproduction strategies.
func (p *SimplePopulation[T]) NextGeneration(selection Selection[T], reproduction Reproduction[T]) *SimplePopulation[T] {
	selection.SetAgents(p.agents)
	return NewSimplePopulation(func() T {
		a, b := selection.Select(), selection.Select()
		if a.Fitness < b.Fitness {
			a, b = b, a
		}
		return reproduction.Reproduce(a.Genotype, b.Genotype)
	}, len(p.agents))
}

// Agents returns the agents in the population.
//
// TODO(change this to an iterator once they get added to the language, as this will increase performance in other cases)
func (p *SimplePopulation[T]) Agents() []*Agent[T] {
	return p.agents
}
