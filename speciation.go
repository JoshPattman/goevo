package goevo

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// Returns a function that computes the genetic distance between the two genotypes.
// This is disjoint * number_of_disjoint_genes + matching * total_matching_synapse_weight_diff.
// The values used in the original paper are disjoint=1, matching=0.4.
// The returned function can be used in Speciate().
func GeneticDistance(disjoint, matching float64) func(*Genotype, *Genotype) float64 {
	return func(g1, g2 *Genotype) float64 {
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
}

// An agent is a genotype with a fitness and a species id.
// It is used when performing the NEAT algorithms.
type Agent struct {
	Genotype  *Genotype
	Fitness   float64
	SpeciesID int
}

// NewAgent creates a new agent with the given genotype.
// It has fitness of 0 and species id of 0.
func NewAgent(g *Genotype) *Agent {
	return &Agent{
		Genotype: g,
	}
}

// Speciate takes a population ([]*Agent) and splits the population into species.
// The distanceThreshold is the maximum distance between two genotypes for them to be considered the same species.
// You must provide a distance function which takes two genotypes and returns a float64 representing the distance between them.
// For instance, you could use the GeneticDistance function.
// If keepExistingSpecies is true, then the species ids of the agents in the population will be preserved.
func Speciate(newSpeciesCounter Counter, agents []*Agent, distanceThreshold float64, keepExistingSpecies bool, distance func(a, b *Genotype) float64) map[int][]*Agent {
	if distanceThreshold < 0 {
		distanceThreshold = 0
	}
	agentsPool := make([]*Agent, len(agents))
	copy(agentsPool, agents)
	species := make(map[int][]*Agent)
	for len(agentsPool) > 0 {
		newSpecies := make([]*Agent, 0)
		agentsNewPool := make([]*Agent, 0)
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
func CalculateOffspring(speciatedPopulation map[int][]*Agent, targetCount int) map[int]int {
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
	speciesAllowedOffspring := make(map[int]int)
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
func Repopulate(speciatedPopulation map[int][]*Agent, allowedOffspringCounts map[int]int, reproduction func(g1, g2 *Genotype) *Genotype, selection func([]*Agent) *Agent) []*Agent {
	// Define new population to fill
	population := make([]*Agent, 0)
	// For every species
	for sid, spec := range speciatedPopulation {
		// Using roulette wheel selection, for the specified number of times, pick two parents proportinate to their fitness
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

// function to perform roulette wheel selection. For this funtion to work, fitnesses must ALWAYS be >= 0.
func RouletteWheelSelection(agents []*Agent) *Agent {
	// calculate total fitness
	totalFitness := 0.0
	for _, agent := range agents {
		fitness := agent.Fitness
		if fitness < 0 {
			panic("Fitness must be > 0")
		}
		totalFitness += fitness
	}
	if totalFitness > 0 {
		// generate random number
		r := rand.Float64() * totalFitness
		// find the index of the selected element
		runningSum := 0.0
		for i, agent := range agents {
			fitness := agent.Fitness
			runningSum += fitness
			if r <= runningSum {
				return agents[i]
			}
		}
		panic("Somthing went wrong with roulette wheel selection")
	} else {
		return agents[rand.Intn(len(agents))]
	}
}

// function to take a slice of agents and return a chosen agent picked probabilistically based on their position when sorted by fitness.
func ProbabilisticSelection(agents []*Agent) *Agent {
	// sort agents by fitness
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].Fitness > agents[j].Fitness
	})
	// generate random number
	r := math.Pow(rand.Float64(), 3)
	// find index of selected agent
	return agents[int(math.Floor(r*float64(len(agents))))]
}
