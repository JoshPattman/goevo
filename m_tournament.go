package goevo

import (
	"math/rand"
)

// tournamentSelection is a tournamentSelection strategy that selects the best agent from a random tournament of agents.
// It implements [tournamentSelection].
type tournamentSelection[T any] struct {
	// The number of agents to include in each tournament.
	tournamentSize int
	agents         []*Agent[T]
}

func NewTournamentSelection[T any](tournamentSize int) Selection[T] {
	if tournamentSize <= 0 {
		panic("must have at least tournament size of 1")
	}
	return &tournamentSelection[T]{
		tournamentSize: tournamentSize,
		agents:         nil,
	}
}

// SetAgents sets the agents to select from for this generation.
func (t *tournamentSelection[T]) SetAgents(agents []*Agent[T]) {
	t.agents = agents
}

// Select returns an agent selected from the population using a tournament.
func (t *tournamentSelection[T]) Select() *Agent[T] {
	if t.agents == nil {
		panic("must call SetAgents before selecting")
	}
	if len(t.agents) == 0 {
		panic("must have at least one agent")
	}
	best := t.agents[rand.Intn(len(t.agents))]
	for i := 0; i < t.tournamentSize-1; i++ {
		testIndex := rand.Intn(len(t.agents))
		if t.agents[testIndex].Fitness > best.Fitness {
			best = t.agents[testIndex]
		}
	}
	return best
}
