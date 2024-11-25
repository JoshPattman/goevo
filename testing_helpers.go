package goevo

import (
	"math"
	"testing"
)

func testWithXORDataset[T any](t *testing.T, pop Population[T], build func(T) Forwarder) {
	X := [][]float64{
		{0, 0, 1},
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
	}
	Y := [][]float64{
		{0},
		{1},
		{1},
		{0},
	}
	testWithDataset(t, X, Y, pop, build)
}

func testWithRecurrentDataset[T any](t *testing.T, pop Population[T], build func(T) Forwarder) {
	X := [][]float64{
		{1},
		{1},
		{1},
		{1},
	}
	Y := [][]float64{
		{0},
		{1},
		{1},
		{0},
	}
	testWithDataset(t, X, Y, pop, build)
}

func testWithDataset[T any](t *testing.T, X, Y [][]float64, pop Population[T], build func(T) Forwarder) {
	if build == nil {
		build = func(t T) Forwarder {
			switch t := any(t).(type) {
			case Buildable:
				return t.Build()
			case Forwarder:
				return t
			default:
				panic("unknown type: cannot build or forward, so you need to supply a build function")
			}
		}

	}
	// We add a bias on the end of each input, which is always 1
	fitness := func(geno T) float64 {
		f := build(geno)
		fitness := 0.0
		for i := range X {
			pred := f.Forward(X[i])
			e := pred[0] - Y[i][0]
			fitness -= math.Pow(math.Abs(e), 3)
		}
		return fitness
	}

	testWithFitnessFunc(t, fitness, pop)
}

func testWithFitnessFunc[T any](t *testing.T, fitness func(T) float64, pop Population[T]) {
	var maxFitness float64
	var maxGt T
	for gen := 0; gen < 5000; gen++ {
		maxFitness = math.Inf(-1)
		maxGt = *new(T)
		for _, a := range pop.All() {
			a.Fitness = fitness(a.Genotype)
			if a.Fitness > maxFitness {
				maxFitness = a.Fitness
				maxGt = a.Genotype
			}
		}

		if maxFitness > -0.1 {
			break
		}

		pop = pop.NextGeneration(nil)
	}
	if maxFitness < -0.1 {
		t.Fatalf("Recurrency Failed to converge, ending with fitness %f", maxFitness)
	}
	val, ok := any(maxGt).(Validateable)
	if ok {
		if err := val.Validate(); err != nil {
			t.Fatalf("final genotype was not valid: %v\nGenotype:\n%v", err, maxGt)
		}
	}
}

func assertEq[T comparable](t *testing.T, a T, b T, name string) {
	if a != b {
		t.Fatalf("error in check '%s' (not equal): '%v' and '%v'", name, a, b)
	}
}
