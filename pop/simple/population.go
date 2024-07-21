// Package simple provides a simple implementation of a population, where all agents are in a single species.
// Agents are selected and bred to create the next generation with the same number of agents.
package simple

import "github.com/JoshPattman/goevo"

var _ goevo.Population[int] = &Population[int]{}

// Population has a single species, and generates the entire next generation by selcting and breeding from the previous one.
type Population[T any] struct {
	Agents       []*goevo.Agent[T]
	Selection    goevo.SelectionStrategy[T]
	Reproduction goevo.Reproduction[T]
}

// NewPopulation creates a new SimplePopulation with n agents, each with a new genotype created by newGenotype.
func NewPopulation[T any](newGenotype func() T, n int, selection goevo.SelectionStrategy[T], reproduction goevo.Reproduction[T]) *Population[T] {
	agents := make([]*goevo.Agent[T], n)
	for i := range agents {
		agents[i] = goevo.NewAgent(newGenotype())
	}
	return &Population[T]{
		Agents:       agents,
		Selection:    selection,
		Reproduction: reproduction,
	}
}

// NextGeneration creates a new SimplePopulation from the current one, using the given selection and reproduction strategies.
func (p *Population[T]) NextGeneration() goevo.Population[T] {
	p.Selection.SetAgents(p.Agents)
	return NewPopulation(func() T {
		parents := goevo.SelectNGenotypes(p.Selection, p.Reproduction.NumParents())
		return p.Reproduction.Reproduce(parents)
	}, len(p.Agents), p.Selection, p.Reproduction)
}

// Agents returns the agents in the population.
//
// TODO(change this to an iterator once they get added to the language, as this will increase performance in other cases)
func (p *Population[T]) All() []*goevo.Agent[T] {
	return p.Agents
}
