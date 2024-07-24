package goevo

import (
	"fmt"
	"math"
	"math/rand"
)

var _ Population[int] = &SpeciatedPopulation[int]{}

// SpeciatedPopulation is a speciated population of agents.
// Each species has the same number of agents, and there are always the same number of species.
// Each generation, with a chance, the worst species is removed, and replaced with a random species or the best species.
type SpeciatedPopulation[T any] struct {
	// Species is a map of species ID to a slice of agents.
	Species map[int][]*Agent[T]
	// RemoveWorstSpeciesChance is the chance that the worst species will be removed each generation.
	RemoveWorstSpeciesChance float64
	// StdNumAgentsSwap is the standard deviation of the number of agents to swap between species each generation.
	// An agent is swapped by moving it to another random species, and moving another agent from that species to this species.
	StdNumAgentsSwap float64
	// Counter is a counter to keep track of the new species.
	Counter *Counter
	// The selection strategy to use when selecting agents to reproduce.
	Selection Selection[T]
	// The reproduction strategy to use when creating new agents.
	Reproduction Reproduction[T]
}

// NewSpeciatedPopulation creates a new speciated population.
func NewSpeciatedPopulation[T any](counter *Counter, newGenotype func() T, numSpecies, numAgentsPerSpecies int, removeWorstSpeciesChance, stdNumAgentsSwap float64, selection Selection[T], reproduction Reproduction[T]) *SpeciatedPopulation[T] {
	species := make(map[int][]*Agent[T])
	for i := 0; i < numSpecies; i++ {
		agents := make([]*Agent[T], numAgentsPerSpecies)
		for j := 0; j < numAgentsPerSpecies; j++ {
			agents[j] = NewAgent(newGenotype())
		}
		species[counter.Next()] = agents
	}
	return &SpeciatedPopulation[T]{
		Species:                  species,
		RemoveWorstSpeciesChance: removeWorstSpeciesChance,
		StdNumAgentsSwap:         stdNumAgentsSwap,
		Counter:                  counter,
		Selection:                selection,
		Reproduction:             reproduction,
	}
}

// NextGeneration implements [Population].
func (p *SpeciatedPopulation[T]) NextGeneration() Population[T] {
	var agentsPerGen int
	for _, agents := range p.Species {
		agentsPerGen = len(agents)
		break
	}
	if agentsPerGen == 0 {
		panic("population: no agents in species, should not have happened")
	}
	numSpecies := len(p.Species)
	// Calculate the average fitness of each species, checking if it is the worst.
	worstFitness := math.Inf(1)
	worstSpecies := 0
	for i, agents := range p.Species {
		var sum float64
		for _, agent := range agents {
			sum += agent.Fitness
		}
		avgFitness := sum / float64(agentsPerGen)
		if avgFitness < worstFitness {
			worstFitness, worstSpecies = avgFitness, i
		}
	}
	// Calculate the ids of the species to reproduce
	type speciesToReproduce struct {
		fromId int
		newId  int
	}
	toReproduce := make([]speciesToReproduce, 0, numSpecies)
	deleteWorst := rand.Float64() < p.RemoveWorstSpeciesChance
	for id := range p.Species {
		if id == worstSpecies && deleteWorst {
			continue
		}
		toReproduce = append(toReproduce, speciesToReproduce{id, id})
	}
	if deleteWorst {
		// pick and index of the species we already know are going to reproduce
		fromIndex := rand.Intn(len(toReproduce))
		// give that species a new id in the next generation
		toReproduce[fromIndex].newId = p.Counter.Next()
		// using the same parent species, add a new species also with a new id
		toReproduce = append(toReproduce, speciesToReproduce{toReproduce[fromIndex].fromId, p.Counter.Next()})
	}
	// Create the new species
	newSpecies := make(map[int][]*Agent[T])
	for _, r := range toReproduce {
		p.Selection.SetAgents(p.Species[r.fromId])
		newAgents := make([]*Agent[T], agentsPerGen)
		for i := range newAgents {
			parents := SelectNGenotypes(p.Selection, p.Reproduction.NumParents())
			newAgents[i] = NewAgent(p.Reproduction.Reproduce(parents))
		}
		newSpecies[r.newId] = newAgents
	}
	// Swap agents between species
	numToSwap := int(math.Round(math.Abs(rand.NormFloat64()) * float64(p.StdNumAgentsSwap)))
	for i := 0; i < numToSwap; i++ {
		aIndex, bIndex := rand.Intn(len(toReproduce)), rand.Intn(len(toReproduce))
		aId, bId := toReproduce[aIndex].newId, toReproduce[bIndex].newId
		aAgentIndex, bAgentIndex := rand.Intn(agentsPerGen), rand.Intn(agentsPerGen)
		newSpecies[aId][aAgentIndex], newSpecies[bId][bAgentIndex] = newSpecies[bId][bAgentIndex], newSpecies[aId][aAgentIndex]
	}
	// Sanity checks
	if len(newSpecies) != numSpecies {
		fmt.Println(len(newSpecies), numSpecies, len(toReproduce), toReproduce)
		panic("population: wrong number of species, should not have happened")
	}
	sanityFoundIds := make(map[int]struct{})
	for id, agents := range newSpecies {
		if len(agents) != agentsPerGen {
			panic("population: wrong number of agents in species, should not have happened")
		}
		if id == worstSpecies && deleteWorst {
			panic("population: worst species was not deleted, should not have happened")
		}
		if _, ok := sanityFoundIds[id]; ok {
			panic("population: duplicate species id, should not have happened")
		}
		sanityFoundIds[id] = struct{}{}
	}
	// Return the new population
	return &SpeciatedPopulation[T]{
		Species:                  newSpecies,
		RemoveWorstSpeciesChance: p.RemoveWorstSpeciesChance,
		StdNumAgentsSwap:         p.StdNumAgentsSwap,
		Counter:                  p.Counter,
		Selection:                p.Selection,
		Reproduction:             p.Reproduction,
	}
}

// All implements [Population].
func (p *SpeciatedPopulation[T]) All() []*Agent[T] {
	all := make([]*Agent[T], 0, len(p.Species)*len(p.Species[0]))
	for _, agents := range p.Species {
		all = append(all, agents...)
	}
	return all
}
