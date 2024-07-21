package goevo

var _ Population[int] = &HillclimberPopulation[int]{}

type HillclimberPopulation[T any] struct {
	A            *Agent[T]
	B            *Agent[T]
	Selection    Selection[T]
	Reproduction Reproduction[T]
}

func NewHillclimberPopulation[T any](initialA, initialB T, selection Selection[T], reproduction Reproduction[T]) *HillclimberPopulation[T] {
	return &HillclimberPopulation[T]{
		A:            NewAgent(initialA),
		B:            NewAgent(initialB),
		Selection:    selection,
		Reproduction: reproduction,
	}
}

func (p *HillclimberPopulation[T]) NextGeneration() Population[T] {
	if p.Reproduction.NumParents() != 1 {
		panic("Hillclimber only supports reproduction with 1 parent")
	}
	p.Selection.SetAgents(p.All())
	parent := p.Selection.Select()
	a := NewAgent(parent.Genotype)
	b := NewAgent(p.Reproduction.Reproduce([]T{parent.Genotype}))
	return &HillclimberPopulation[T]{A: a, B: b}
}

func (p *HillclimberPopulation[T]) All() []*Agent[T] {
	return []*Agent[T]{p.A, p.B}
}
