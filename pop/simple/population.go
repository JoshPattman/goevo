// Package simple provides a simple implementation of a population, where all agents are in a single species.
// Agents are selected and bred to create the next generation with the same number of agents.
package simple

import "github.com/JoshPattman/goevo"

// Population has a single species, and generates the entire next generation by selcting and breeding from the previous one.
type Population[T any] struct {
	agents []*goevo.Agent[T]
}

// NewPopulation creates a new SimplePopulation with n agents, each with a new genotype created by newGenotype.
func NewPopulation[T any](newGenotype func() T, n int) *Population[T] {
	agents := make([]*goevo.Agent[T], n)
	for i := range agents {
		agents[i] = goevo.NewAgent(newGenotype())
	}
	return &Population[T]{
		agents: agents,
	}
}

// NextGeneration creates a new SimplePopulation from the current one, using the given selection and reproduction strategies.
func (p *Population[T]) NextGeneration(selection goevo.Selection[T], reproduction goevo.Reproduction[T]) *Population[T] {
	selection.SetAgents(p.agents)
	return NewPopulation(func() T {
		parents := make([]T, reproduction.NumParents())
		for i := range parents {
			parents[i] = selection.Select().Genotype
		}
		return reproduction.Reproduce(parents)
	}, len(p.agents))
}

// Agents returns the agents in the population.
//
// TODO(change this to an iterator once they get added to the language, as this will increase performance in other cases)
func (p *Population[T]) Agents() []*goevo.Agent[T] {
	return p.agents
}
