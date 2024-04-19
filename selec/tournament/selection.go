// Package tournament provides an implementation of the Selection interface that selects the best agent from a random tournament of n agents.
package tournament

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

// Ensure that TournamentSelection implements Selection.
var _ goevo.Selection[int] = &Selection[int]{}

// Selection is a selection strategy that selects the best agent from a random tournament of agents.
// It implements [Selection].
type Selection[T any] struct {
	// The number of agents to include in each tournament.
	TournamentSize int
	agents         []*goevo.Agent[T]
}

// SetAgents sets the agents to select from for this generation.
func (t *Selection[T]) SetAgents(agents []*goevo.Agent[T]) {
	t.agents = agents
}

// Select returns an agent selected from the population using a tournament.
func (t *Selection[T]) Select() *goevo.Agent[T] {
	if t.agents == nil {
		panic("must call SetAgents before selecting")
	}
	if len(t.agents) == 0 {
		panic("must have at least one agent")
	}
	if t.TournamentSize <= 0 {
		panic("must have tournamnet size of at least 1")
	}
	best := t.agents[rand.Intn(len(t.agents))]
	for i := 0; i < t.TournamentSize-1; i++ {
		testIndex := rand.Intn(len(t.agents))
		if t.agents[testIndex].Fitness > best.Fitness {
			best = t.agents[testIndex]
		}
	}
	return best
}
