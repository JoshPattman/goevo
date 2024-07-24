package goevo

import (
	"math"
	"testing"
)

func TestFloatsGt(t *testing.T) {
	counter := NewCounter()
	mut := &ArrayMutationStd[float64]{
		MutateProbability: 0.1,
		MutateStd:         0.05,
	}
	crs := &ArrayCrossoverKPoint[float64]{K: 2}
	reprod := NewTwoPhaseReproduction(crs, mut)
	selec := &TournamentSelection[*ArrayGenotype[float64]]{
		TournamentSize: 3,
	}
	var pop Population[*ArrayGenotype[float64]] = NewSpeciatedPopulation(
		counter,
		func() *ArrayGenotype[float64] { return NewFloatArrayGenotype(10, 0.5) },
		5,
		20,
		0.1,
		2.5,
		selec,
		reprod,
	)
	// Fitness is max (0) when all the numbers sum to 10
	fitness := func(f *ArrayGenotype[float64]) float64 {
		total := 0.0
		for i := range f.Values {
			total += f.Values[i]
		}
		return -math.Abs(10 - total)
	}
	var highestFitness float64
	for gen := 0; gen < 100; gen++ {
		highestFitness = math.Inf(-1)
		for _, a := range pop.All() {
			a.Fitness = fitness(a.Genotype)
			if a.Fitness > highestFitness {
				highestFitness = a.Fitness
			}
		}
		pop = pop.NextGeneration()
	}
	if highestFitness < -0.1 {
		t.Fatalf("Failed to converge, ending with fitness %f", highestFitness)
	}
}
