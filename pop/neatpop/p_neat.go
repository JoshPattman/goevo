package goevo

import "github.com/JoshPattman/goevo"

// Population is a population of agents, each with a genotype of type T.
// The population is split into species, and the next generation is created by selecting and breeding from the previous one,
// while respecting the species.
// It uses the Nero Evolution of Augmenting Topologies (NEAT) algorithm.
type Population[T any] struct {
	// The maximum distance between two genotypes for them to be considered the same species.
	DistanceThreshold float64
	// The multiplier for the distance threshold, which is used to increase or decrease the distance threshold each generation.
	DistanceThresholdMultiplier float64
	species                     [][]*goevo.Agent[T]
}

// NewPopulation creates a new NEATPopulation with n agents, each with a new genotype created by newGenotype.
func NewPopulation[T any](newGenotype func() T, n int, initialDistanceThreshold, distanceThresholdMultiplier float64) *Population[T] {
	pop := &Population[T]{
		species:                     make([][]*goevo.Agent[T], 1),
		DistanceThreshold:           initialDistanceThreshold,
		DistanceThresholdMultiplier: distanceThresholdMultiplier,
	}
	pop.species[0] = make([]*goevo.Agent[T], n)
	for i := range pop.species[0] {
		pop.species[0][i] = goevo.NewAgent(newGenotype())
	}
	return pop
}

/* WIP
import (
	"fmt"
	"math"
	"math/rand"
)

// NEATPopulation is a slice of agents
type NEATPopulation []*Agent

// SpeciatedPopulation is a population, but split into species
type SpeciatedPopulation map[int][]*Agent

// OffspringCounts stores how many offspring each species can have
type OffspringCounts map[int]int

// This is disjoint * number_of_disjoint_genes + matching * total_matching_synapse_weight_diff.
// The values used in the original paper are disjoint=1, matching=0.4.
func geneticDistance(disjoint, matching float64, g1, g2 *Genotype) float64 {
	numMatching := 0.0
	totalWDiff := 0.0
	for sid1, s1 := range g1.Synapses {
		if s2, ok := g2.Synapses[sid1]; ok {
			// Matching
			numMatching++
			totalWDiff += math.Abs(s1.Weight - s2.Weight)
		}
	}
	numDisjoint := (float64(len(g1.Synapses)) - numMatching) + (float64(len(g2.Synapses)) - numMatching)

	return disjoint*numDisjoint + matching*totalWDiff
}

// Speciate takes a population ([]*Agent) and splits the population into species.
// The distanceThreshold is the maximum distance between two genotypes for them to be considered the same species.
// You must provide a distance function which takes two genotypes and returns a float64 representing the distance between them.
// For instance, you could use the GeneticDistance function.
// If keepExistingSpecies is true, then the species ids of the agents in the population will be preserved.
func (population NEATPopulation) Speciate(newSpeciesCounter *Counter, distanceThreshold float64, keepExistingSpecies bool, distance GeneticDistanceFunc) SpeciatedPopulation {
	if distanceThreshold < 0 {
		distanceThreshold = 0
	}
	agentsPool := make(NEATPopulation, len(population))
	copy(agentsPool, population)
	species := make(SpeciatedPopulation)
	for len(agentsPool) > 0 {
		newSpecies := make(NEATPopulation, 0)
		agentsNewPool := make(NEATPopulation, 0)
		template := agentsPool[rand.Intn(len(agentsPool))]
		for _, a := range agentsPool {
			if distance(template.Genotype, a.Genotype) <= distanceThreshold {
				// Same species
				newSpecies = append(newSpecies, a)
			} else {
				agentsNewPool = append(agentsNewPool, a)
			}
		}
		agentsPool = agentsNewPool
		if keepExistingSpecies {
			existingSpecies, speciesExists := species[template.SpeciesID]
			if speciesExists {
				// Preserve the larger of the two species and assign a new species id for the smaller one
				newID := newSpeciesCounter.Next()
				if len(existingSpecies) > len(newSpecies) {
					species[newID] = newSpecies
				} else {
					species[newID] = species[template.SpeciesID]
					species[template.SpeciesID] = newSpecies
				}
			} else {
				species[template.SpeciesID] = newSpecies
			}
		} else {
			newID := newSpeciesCounter.Next()
			species[newID] = newSpecies
		}
	}
	for sid := range species {
		for _, a := range species[sid] {
			a.SpeciesID = sid
		}
	}
	return species
}

// Calculate the number of offspring each species should have.
// targetCount is the total number of offspring to be created (the actual number of offspring may vary slightly due to rounding).
func (speciatedPopulation SpeciatedPopulation) CalculateOffspring(targetCount int) OffspringCounts {
	// Caclulate the total fitness of each species and the global total fitness (adjusted fitness)
	speciesTotalFitness := make(map[int]float64)
	globalTotalFitness := 0.0
	for sid, spec := range speciatedPopulation {
		totalFitness := 0.0
		for _, agent := range spec {
			if agent.Fitness < 0 {
				panic("Fitness is less than 0. This cannot happen")
			}
			totalFitness += agent.Fitness / float64(len(spec))
		}
		speciesTotalFitness[sid] = totalFitness
		globalTotalFitness += totalFitness
	}
	// Calculate the number of offspring each species should have
	speciesAllowedOffspring := make(OffspringCounts)
	for sid := range speciatedPopulation {
		speciesAllowedOffspring[sid] = int(math.Round(float64(targetCount) * speciesTotalFitness[sid] / globalTotalFitness))
		if speciesAllowedOffspring[sid] < 0 {
			fmt.Println(targetCount, speciesTotalFitness, globalTotalFitness)
			fmt.Println(speciatedPopulation)
			fmt.Println("")
			for _, s := range speciatedPopulation {
				for _, a := range s {
					fmt.Print(a.Fitness, " ")
				}
				fmt.Println("")
			}
			panic("Number of offspring is less than 0. This cannot happen")
		}
	}
	return speciesAllowedOffspring
}

// Create a new population by picking parents from each species and creating a child from them.
// Fitnesses must be assigned to the agents before calling this function.
// allowedOffspringCounts is a map of species ids to the number of offspring that species is allowed to have, which can be obtained by using CalculateOffspring.
// reproduction is a function which takes two genotypes and returns a new genotype which is the child of the two parents.
// selection is a function which takes a slice of agents and returns a single agent which is the parent.
func (speciatedPopulation SpeciatedPopulation) NextGeneration(allowedOffspringCounts OffspringCounts, reproduction ReproductionFunc, selection SelectionFunc) NEATPopulation {
	// Define new population to fill
	population := make(NEATPopulation, 0)
	// For every species
	for sid, spec := range speciatedPopulation {
		// Using selection, for the specified number of times, pick two parents proportinate to their fitness
		// Create a new agent which is the child of both parents. Ensure the first parent is the more fit one
		// Add the new agent to the new population
		for i := 0; i < allowedOffspringCounts[sid]; i++ {
			// Pick two parents
			parent1 := selection(spec)
			parent2 := selection(spec)
			// Ensure parent1 is the more fit one
			if parent1.Fitness < parent2.Fitness {
				parent1, parent2 = parent2, parent1
			}
			// Create a new agent
			child := NewAgent(reproduction(parent1.Genotype, parent2.Genotype))
			child.SpeciesID = sid
			// Add the new agent to the new population
			population = append(population, child)
		}
	}
	return population
}

func (s SpeciatedPopulation) GetBestByMean() (int, float64) {
	maxAvgFitness := math.Inf(-1)
	maxSpecies := -1
	for si := range s {
		totalFitness := 0.0
		for _, a := range s[si] {
			totalFitness += a.Fitness
		}
		avgFitness := totalFitness / float64(len(s[si]))
		if avgFitness > maxAvgFitness {
			maxAvgFitness = avgFitness
			maxSpecies = si
		}
	}
	return maxSpecies, maxAvgFitness
}

func (s SpeciatedPopulation) GetWorstByMean() (int, float64) {
	minAvgFitness := math.Inf(1)
	minSpecies := -1
	for si := range s {
		totalFitness := 0.0
		for _, a := range s[si] {
			totalFitness += a.Fitness
		}
		avgFitness := totalFitness / float64(len(s[si]))
		if avgFitness < minAvgFitness {
			minAvgFitness = avgFitness
			minSpecies = si
		}
	}
	return minSpecies, minAvgFitness
}

func (s SpeciatedPopulation) GetBestByMax() (int, float64) {
	maxFitness := math.Inf(-1)
	maxSpecies := -1
	for si := range s {
		for _, a := range s[si] {
			if a.Fitness > maxFitness {
				maxFitness = a.Fitness
				maxSpecies = si
			}
		}
	}
	return maxSpecies, maxFitness
}

func (s SpeciatedPopulation) GetWorstByMax() (int, float64) {
	minFitness := math.Inf(1)
	minSpecies := -1
	for si := range s {
		speciesMax := math.Inf(-1)
		for _, a := range s[si] {
			if a.Fitness > speciesMax {
				speciesMax = a.Fitness
			}
		}
		if speciesMax < minFitness {
			minFitness = speciesMax
			minSpecies = si
		}
	}
	return minSpecies, minFitness
}
*/
