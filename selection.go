package goevo

import "math/rand"

type Selection interface {
	SetAgents(agents []*Agent)
	Select() *Agent
}

type TournamentSelection struct {
	TournamentSize int
	agents         []*Agent
}

func (t *TournamentSelection) SetAgents(agents []*Agent) {
	t.agents = agents
}

func (t *TournamentSelection) Select() *Agent {
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
