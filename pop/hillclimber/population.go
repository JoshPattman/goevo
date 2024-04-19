// Package hillclimber provides an implementation of the Population interface that maintains two agents and selects the best of the two for the next generation.
package hillclimber

import "github.com/JoshPattman/goevo"

type Population[T any] struct {
	A *goevo.Agent[T]
	B *goevo.Agent[T]
}

func NewPopulation[T any](initialA, initialB T) *Population[T] {
	return &Population[T]{A: goevo.NewAgent(initialA), B: goevo.NewAgent(initialB)}
}

func (p *Population[T]) NextGeneration(selection goevo.Selection[T], reproduction goevo.Reproduction[T]) *Population[T] {
	if reproduction.NumParents() != 1 {
		panic("Hillclimber only supports reproduction with 1 parent")
	}
	selection.SetAgents(p.Agents())
	parent := selection.Select()
	a := goevo.NewAgent(parent.Genotype)
	b := goevo.NewAgent(reproduction.Reproduce([]T{parent.Genotype}))
	return &Population[T]{A: a, B: b}
}

func (p *Population[T]) Agents() []*goevo.Agent[T] {
	return []*goevo.Agent[T]{p.A, p.B}
}
