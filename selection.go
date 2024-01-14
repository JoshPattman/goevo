package goevo

import (
	"math"
	"math/rand"
	"sort"
)

type SelectionFunc func([]*Agent) *Agent

// function to perform roulette wheel selection. For this funtion to work, fitnesses must ALWAYS be >= 0.
func RouletteWheelSecelction() SelectionFunc {
	return func(agents []*Agent) *Agent {
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
}

// Degree of 3 is good
func PolyProbSelection(polyDegree float64) SelectionFunc {
	return func(agents []*Agent) *Agent {
		// sort agents by fitness
		sort.Slice(agents, func(i, j int) bool {
			return agents[i].Fitness > agents[j].Fitness
		})
		// generate random number
		r := math.Pow(rand.Float64(), polyDegree)
		// find index of selected agent
		return agents[int(math.Floor(r*float64(len(agents))))]
	}
}
