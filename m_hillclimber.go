package goevo

var _ Population[int] = &HillClimberPopulation[int]{}

type HillClimberPopulation[T any] struct {
	A            *Agent[T]
	B            *Agent[T]
	Selection    Selection[T]
	Reproduction Reproduction[T]
}

func NewHillClimberPopulation[T any](initialA, initialB T, selection Selection[T], reproduction Reproduction[T]) *HillClimberPopulation[T] {
	return &HillClimberPopulation[T]{
		A:            NewAgent(initialA),
		B:            NewAgent(initialB),
		Selection:    selection,
		Reproduction: reproduction,
	}
}

func (p *HillClimberPopulation[T]) NextGeneration() Population[T] {
	if p.Reproduction.NumParents() != 1 {
		panic("Hillclimber only supports reproduction with 1 parent")
	}
	p.Selection.SetAgents(p.All())
	parent := p.Selection.Select()
	a := NewAgent(parent.Genotype)
	b := NewAgent(p.Reproduction.Reproduce([]T{parent.Genotype}))
	return &HillClimberPopulation[T]{A: a, B: b}
}

func (p *HillClimberPopulation[T]) All() []*Agent[T] {
	return []*Agent[T]{p.A, p.B}
}
