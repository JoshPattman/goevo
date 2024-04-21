// Package elite provides an implementation of the Selection interface that selects the best agent from the previous generation.
package elite

import (
	"math"

	"github.com/JoshPattman/goevo"
)

// Ensure that Selection implements Selection.
var _ goevo.SelectionStrategy[int] = &Selection[int]{}

type Selection[T any] struct {
	lastBest *goevo.Agent[T]
}

func (s *Selection[T]) SetAgents(agents []*goevo.Agent[T]) {
	bestFitness := math.Inf(-1)
	var bestAgent *goevo.Agent[T]
	for _, agent := range agents {
		if agent.Fitness > bestFitness {
			bestFitness = agent.Fitness
			bestAgent = agent
		}
	}
	s.lastBest = bestAgent
}

func (s *Selection[T]) Select() *goevo.Agent[T] {
	if s.lastBest == nil {
		panic("must call SetAgents before selecting (also ensure at least one agent has fitness greater than -inf)")
	}
	return s.lastBest
}
