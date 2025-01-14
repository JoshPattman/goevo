package goevo

import (
	"math/rand"
)

// TournamentSelection is a TournamentSelection strategy that selects the best agent from a random tournament of agents.
// It implements [TournamentSelection].
type TournamentSelection[T any] struct {
	// The number of agents to include in each tournament.
	TournamentSize int
	agents         []*Agent[T]
}

// SetAgents sets the agents to select from for this generation.
func (t *TournamentSelection[T]) SetAgents(agents []*Agent[T]) {
	t.agents = agents
}

// Select returns an agent selected from the population using a tournament.
func (t *TournamentSelection[T]) Select() *Agent[T] {
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
