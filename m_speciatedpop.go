package goevo

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
)

// SpeciatedPopulation is a speciated population of agents.
// Each species has the same number of agents, and there are always the same number of species.
// Each generation, with a chance, the worst species is removed, and replaced with a random species or the best species.
type SpeciatedPopulation[T any] struct {
	// species is a map of species ID to a slice of agents.
	species map[int][]*Agent[T]
	// removeWorstSpeciesChance is the chance that the worst species will be removed each generation.
	removeWorstSpeciesChance float64
	// stdNumAgentsSwap is the standard deviation of the number of agents to swap between species each generation.
	// An agent is swapped by moving it to another random species, and moving another agent from that species to this species.
	stdNumAgentsSwap float64
	// counter is a counter to keep track of the new species.
	counter *Counter
	// The selection strategy to use when selecting agents to reproduce.
	selection Selection[T]
	// The reproduction strategy to use when creating new agents.
	reproduction Reproduction[T]
}

// NewSpeciatedPopulation creates a new speciated population.
func NewSpeciatedPopulation[T any](
	counter *Counter,
	newGenotype func() T,
	numSpecies int,
	numAgentsPerSpecies int,
	removeWorstSpeciesChance float64,
	stdNumAgentsSwap float64,
	selection Selection[T],
	reproduction Reproduction[T],
) *SpeciatedPopulation[T] {
	species := make(map[int][]*Agent[T])
	for i := 0; i < numSpecies; i++ {
		agents := make([]*Agent[T], numAgentsPerSpecies)
		for j := 0; j < numAgentsPerSpecies; j++ {
			agents[j] = NewAgent(newGenotype())
		}
		species[counter.Next()] = agents
	}
	return NewSpeciatedPopulationFrom(counter, species, removeWorstSpeciesChance, stdNumAgentsSwap, selection, reproduction)
}

func NewSpeciatedPopulationFrom[T any](
	counter *Counter,
	species map[int][]*Agent[T],
	removeWorstSpeciesChance float64,
	stdNumAgentsSwap float64,
	selection Selection[T],
	reproduction Reproduction[T],
) *SpeciatedPopulation[T] {
	if len(species) == 0 {
		panic("can only create speciated population with at least one species")
	}
	if selection == nil {
		panic("cannot have nil selection")
	}
	if reproduction == nil {
		panic("cannot have nil reproduction")
	}
	firstKey := -1
	for k := range species {
		firstKey = k
		break
	}
	speciesLength := len(species[firstKey])
	if speciesLength == 0 {
		panic("species must have at least one member")
	}
	speciesCopy := make(map[int][]*Agent[T])
	for id, agents := range species {
		if len(agents) != speciesLength {
			panic("all species must be same length")
		}
		speciesCopy[id] = slices.Clone(agents)
	}
	return &SpeciatedPopulation[T]{
		species:                  speciesCopy,
		removeWorstSpeciesChance: removeWorstSpeciesChance,
		stdNumAgentsSwap:         stdNumAgentsSwap,
		counter:                  counter,
		selection:                selection,
		reproduction:             reproduction,
	}
}

// NextGeneration implements [Population].
func (p *SpeciatedPopulation[T]) NextGeneration() Population[T] {
	var agentsPerGen int
	for _, agents := range p.species {
		agentsPerGen = len(agents)
		break
	}
	if agentsPerGen == 0 {
		panic("population: no agents in species, should not have happened")
	}
	numSpecies := len(p.species)
	// Calculate the average fitness of each species, checking if it is the worst.
	worstFitness := math.Inf(1)
	worstSpecies := 0
	for i, agents := range p.species {
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
	deleteWorst := rand.Float64() < p.removeWorstSpeciesChance
	for id := range p.species {
		if id == worstSpecies && deleteWorst {
			continue
		}
		toReproduce = append(toReproduce, speciesToReproduce{id, id})
	}
	if deleteWorst {
		// pick and index of the species we already know are going to reproduce
		fromIndex := rand.Intn(len(toReproduce))
		// give that species a new id in the next generation
		toReproduce[fromIndex].newId = p.counter.Next()
		// using the same parent species, add a new species also with a new id
		toReproduce = append(toReproduce, speciesToReproduce{toReproduce[fromIndex].fromId, p.counter.Next()})
	}
	// Create the new species
	newSpecies := make(map[int][]*Agent[T])
	for _, r := range toReproduce {
		p.selection.SetAgents(p.species[r.fromId])
		newAgents := make([]*Agent[T], agentsPerGen)
		for i := range newAgents {
			parents := SelectNGenotypes(p.selection, p.reproduction.NumParents())
			newAgents[i] = NewAgent(p.reproduction.Reproduce(parents))
		}
		newSpecies[r.newId] = newAgents
	}
	// Swap agents between species
	numToSwap := int(math.Round(math.Abs(rand.NormFloat64()) * float64(p.stdNumAgentsSwap)))
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
		species:                  newSpecies,
		removeWorstSpeciesChance: p.removeWorstSpeciesChance,
		stdNumAgentsSwap:         p.stdNumAgentsSwap,
		counter:                  p.counter,
		selection:                p.selection,
		reproduction:             p.reproduction,
	}
}

// All implements [Population].
func (p *SpeciatedPopulation[T]) All() []*Agent[T] {
	all := make([]*Agent[T], 0, len(p.species)*len(p.species[0]))
	for _, agents := range p.species {
		all = append(all, agents...)
	}
	return all
}

func (p *SpeciatedPopulation[T]) AllSpecies() map[int][]*Agent[T] {
	res := make(map[int][]*Agent[T])
	for id, agents := range p.species {
		res[id] = slices.Clone(agents)
	}
	return res
}
