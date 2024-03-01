package goevo

import "math/rand"

// Selection is a strategy for selecting agents from a population.
// It acts on agents of type T.
type Selection[T any] interface {
	// SetAgents sets the agents to select from for this generation.
	SetAgents(agents []*Agent[T])
	// Select returns an agent selected from the population.
	Select() *Agent[T]
}

// Ensure that TournamentSelection implements Selection.
var _ Selection[*NEATGenotype] = &TournamentSelection[*NEATGenotype]{}

// TournamentSelection is a selection strategy that selects the best agent from a random tournament of agents.
// It implements [Selection].
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
