// Package hillclimber provides an implementation of the Population interface that maintains two agents and selects the best of the two for the next generation.
package hillclimber

import "github.com/JoshPattman/goevo"

var _ goevo.Population[int] = &Population[int]{}

type Population[T any] struct {
	A            *goevo.Agent[T]
	B            *goevo.Agent[T]
	Selection    goevo.SelectionStrategy[T]
	Reproduction goevo.ReproductionStrategy[T]
}

func NewPopulation[T any](initialA, initialB T, selection goevo.SelectionStrategy[T], reproduction goevo.ReproductionStrategy[T]) *Population[T] {
	return &Population[T]{
		A:            goevo.NewAgent(initialA),
		B:            goevo.NewAgent(initialB),
		Selection:    selection,
		Reproduction: reproduction,
	}
}

func (p *Population[T]) NextGeneration() goevo.Population[T] {
	if p.Reproduction.NumParents() != 1 {
		panic("Hillclimber only supports reproduction with 1 parent")
	}
	p.Selection.SetAgents(p.All())
	parent := p.Selection.Select()
	a := goevo.NewAgent(parent.Genotype)
	b := goevo.NewAgent(p.Reproduction.Reproduce([]T{parent.Genotype}))
	return &Population[T]{A: a, B: b}
}

func (p *Population[T]) All() []*goevo.Agent[T] {
	return []*goevo.Agent[T]{p.A, p.B}
}
