package goevo

type HillClimberPopulation[T any] struct {
	a            *Agent[T]
	b            *Agent[T]
	selection    Selection[T]
	reproduction Reproduction[T]
}

func NewHillClimberPopulation[T any](initialA, initialB T, selection Selection[T], reproduction Reproduction[T]) *HillClimberPopulation[T] {
	if selection == nil {
		panic("cannot have nil selection")
	}
	if reproduction == nil {
		panic("cannot have nil reproduction")
	}
	return &HillClimberPopulation[T]{
		a:            NewAgent(initialA),
		b:            NewAgent(initialB),
		selection:    selection,
		reproduction: reproduction,
	}
}

func (p *HillClimberPopulation[T]) NextGeneration() Population[T] {
	if p.reproduction.NumParents() != 1 {
		panic("Hillclimber only supports reproduction with 1 parent")
	}
	p.selection.SetAgents(p.All())
	parent := p.selection.Select()
	a := NewAgent(parent.Genotype)
	b := NewAgent(p.reproduction.Reproduce([]T{parent.Genotype}))
	return &HillClimberPopulation[T]{a: a, b: b, selection: p.selection, reproduction: p.reproduction}
}

func (p *HillClimberPopulation[T]) All() []*Agent[T] {
	return []*Agent[T]{p.a, p.b}
}

func (p *HillClimberPopulation[T]) Both() (*Agent[T], *Agent[T]) {
	return p.a, p.b
}
