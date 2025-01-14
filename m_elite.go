package goevo

import "math"

type EliteSelection[T any] struct {
	lastBest *Agent[T]
}

func (s *EliteSelection[T]) SetAgents(agents []*Agent[T]) {
	bestFitness := math.Inf(-1)
	var bestAgent *Agent[T]
	for _, agent := range agents {
		if agent.Fitness > bestFitness {
			bestFitness = agent.Fitness
			bestAgent = agent
		}
	}
	s.lastBest = bestAgent
}

func (s *EliteSelection[T]) Select() *Agent[T] {
	if s.lastBest == nil {
		panic("must call SetAgents before selecting (also ensure at least one agent has fitness greater than -inf)")
	}
	return s.lastBest
}
