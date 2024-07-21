// Package simple provides a simple implementation of a population, where all agents are in a single species.
// Agents are selected and bred to create the next generation with the same number of agents.
package simple

import "github.com/JoshPattman/goevo"

var _ goevo.Population[int] = &SimplePopulation[int]{}

// SimplePopulation has a single species, and generates the entire next generation by selcting and breeding from the previous one.
type SimplePopulation[T any] struct {
	Agents       []*goevo.Agent[T]
	Selection    goevo.Selection[T]
	Reproduction goevo.Reproduction[T]
}

// NewSimplePopulation creates a new SimplePopulation with n agents, each with a new genotype created by newGenotype.
func NewSimplePopulation[T any](newGenotype func() T, n int, selection goevo.Selection[T], reproduction goevo.Reproduction[T]) *SimplePopulation[T] {
	agents := make([]*goevo.Agent[T], n)
	for i := range agents {
		agents[i] = goevo.NewAgent(newGenotype())
	}
	return &SimplePopulation[T]{
		Agents:       agents,
		Selection:    selection,
		Reproduction: reproduction,
	}
}

// NextGeneration creates a new SimplePopulation from the current one, using the given selection and reproduction strategies.
func (p *SimplePopulation[T]) NextGeneration() goevo.Population[T] {
	p.Selection.SetAgents(p.Agents)
	return NewSimplePopulation(func() T {
		parents := goevo.SelectNGenotypes(p.Selection, p.Reproduction.NumParents())
		return p.Reproduction.Reproduce(parents)
	}, len(p.Agents), p.Selection, p.Reproduction)
}

// Agents returns the agents in the population.
//
// TODO(change this to an iterator once they get added to the language, as this will increase performance in other cases)
func (p *SimplePopulation[T]) All() []*goevo.Agent[T] {
	return p.Agents
}
