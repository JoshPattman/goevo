package goevo

// SimplePopulation has a single species, and generates the entire next generation by selcting and breeding from the previous
type SimplePopulation struct {
	agents []*Agent
}

func NewSimplePopulation(newGenotype func() *Genotype, n int) *SimplePopulation {
	agents := make([]*Agent, n)
	for i := range agents {
		agents[i] = NewAgent(newGenotype())
	}
	return &SimplePopulation{
		agents: agents,
	}
}

func (p *SimplePopulation) NextGeneration(selection Selection, reproduction Reproduction) *SimplePopulation {
	selection.SetAgents(p.agents)
	return NewSimplePopulation(func() *Genotype {
		a, b := selection.Select(), selection.Select()
		if a.Fitness < b.Fitness {
			a, b = b, a
		}
		return reproduction.Reproduce(a.Genotype, b.Genotype)
	}, len(p.agents))
}

func (p *SimplePopulation) Agents() []*Agent {
	return p.agents
}
