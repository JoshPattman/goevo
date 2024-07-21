// Package hillclimber provides an implementation of the Population interface that maintains two agents and selects the best of the two for the next generation.
package hillclimber

import "github.com/JoshPattman/goevo"

var _ goevo.Population[int] = &HillclimberPopulation[int]{}

type HillclimberPopulation[T any] struct {
	A            *goevo.Agent[T]
	B            *goevo.Agent[T]
	Selection    goevo.Selection[T]
	Reproduction goevo.Reproduction[T]
}

func NewHillclimberPopulation[T any](initialA, initialB T, selection goevo.Selection[T], reproduction goevo.Reproduction[T]) *HillclimberPopulation[T] {
	return &HillclimberPopulation[T]{
		A:            goevo.NewAgent(initialA),
		B:            goevo.NewAgent(initialB),
		Selection:    selection,
		Reproduction: reproduction,
	}
}

func (p *HillclimberPopulation[T]) NextGeneration() goevo.Population[T] {
	if p.Reproduction.NumParents() != 1 {
		panic("Hillclimber only supports reproduction with 1 parent")
	}
	p.Selection.SetAgents(p.All())
	parent := p.Selection.Select()
	a := goevo.NewAgent(parent.Genotype)
	b := goevo.NewAgent(p.Reproduction.Reproduce([]T{parent.Genotype}))
	return &HillclimberPopulation[T]{A: a, B: b}
}

func (p *HillclimberPopulation[T]) All() []*goevo.Agent[T] {
	return []*goevo.Agent[T]{p.A, p.B}
}
