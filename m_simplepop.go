package goevo

// SimplePopulation has a single species, and generates the entire next generation by selcting and breeding from the previous one.
type SimplePopulation[T any] struct {
	agents       []*Agent[T]
	selection    Selection[T]
	reproduction Reproduction[T]
}

// NewSimplePopulation creates a new SimplePopulation with n agents, each with a new genotype created by newGenotype.
func NewSimplePopulation[T any](newGenotype func() T, n int, selection Selection[T], reproduction Reproduction[T]) *SimplePopulation[T] {
	if selection == nil {
		panic("cannot have nil selection")
	}
	if reproduction == nil {
		panic("cannot have nil reproduction")
	}
	if n <= 0 {
		panic("cannot create population with less than 1 member")
	}
	agents := make([]*Agent[T], n)
	for i := range agents {
		agents[i] = NewAgent(newGenotype())
	}
	return &SimplePopulation[T]{
		agents:       agents,
		selection:    selection,
		reproduction: reproduction,
	}
}

// NextGeneration creates a new SimplePopulation from the current one, using the given selection and reproduction strategies.
func (p *SimplePopulation[T]) NextGeneration() Population[T] {
	p.selection.SetAgents(p.agents)
	return NewSimplePopulation(func() T {
		parents := SelectNGenotypes(p.selection, p.reproduction.NumParents())
		return p.reproduction.Reproduce(parents)
	}, len(p.agents), p.selection, p.reproduction)
}

// Agents returns the agents in the population.
//
// TODO(change this to an iterator once they get added to the language, as this will increase performance in other cases)
func (p *SimplePopulation[T]) All() []*Agent[T] {
	return p.agents
}
