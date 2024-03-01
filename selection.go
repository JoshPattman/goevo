package goevo

import "math/rand"

type Selection[T any] interface {
	SetAgents(agents []*Agent[T])
	Select() *Agent[T]
}

var _ Selection[*NEATGenotype] = &TournamentSelection[*NEATGenotype]{}

type TournamentSelection[T any] struct {
	TournamentSize int
	agents         []*Agent[T]
}

func (t *TournamentSelection[T]) SetAgents(agents []*Agent[T]) {
	t.agents = agents
}

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
