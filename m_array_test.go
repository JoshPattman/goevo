package goevo

import (
	"math"
	"testing"
)

func setupArrayTestStuff[T any](mut Mutation[*ArrayGenotype[T]], newGenotype func() *ArrayGenotype[T], crsType int, selecType int) Population[*ArrayGenotype[T]] {
	counter := NewCounter()
	var crs Crossover[*ArrayGenotype[T]]
	switch crsType {
	case 0:
		crs = NewArrayCrossoverKPoint[T](2)
	case 1:
		crs = NewArrayCrossoverUniform[T]()
	case 2:
		crs = NewArrayCrossoverAsexual[T]()
	}
	reprod := NewTwoPhaseReproduction(crs, mut)
	var selec Selection[*ArrayGenotype[T]]
	switch selecType {
	case 0:
		selec = NewTournamentSelection[*ArrayGenotype[T]](3)
	case 1:
		selec = &eliteSelection[*ArrayGenotype[T]]{}
	}
	var pop Population[*ArrayGenotype[T]]
	if crsType == 2 {
		g1 := newGenotype()
		pop = NewHillClimberPopulation(
			g1,
			Clone(g1),
			selec,
			reprod,
		)
	} else {
		pop = NewSpeciatedPopulation(
			counter,
			newGenotype,
			5,
			20,
			0.1,
			2.5,
			selec,
			reprod,
		)
	}
	return pop
}

func TestArrayGenotype(t *testing.T) {
	mut := NewArrayMutationGeneratorAdd(NewGeneratorNormal(0, 0.05), 0.1)
	newGenotype := func() *ArrayGenotype[float64] {
		return NewArrayGenotype(10, NewGeneratorNormal(0, 0.5))
	}
	pop := setupArrayTestStuff(mut, newGenotype, 0, 0)
	// Fitness is max (0) when all the numbers sum to 10
	fitness := func(f *ArrayGenotype[float64]) float64 {
		total := 0.0
		for i := range f.values {
			total += f.values[i]
		}
		return -math.Abs(10 - total)
	}
	testWithFitnessFunc(t, fitness, pop)
}

func TestRuneGenotype(t *testing.T) {
	runeset := []rune("ab")
	valueGen := NewGeneratorChoices(runeset)
	mut := NewArrayMutationGeneratorReplace(valueGen, 0.1)
	newGenotype := func() *ArrayGenotype[rune] { return NewArrayGenotype(10, valueGen) }
	pop := setupArrayTestStuff(mut, newGenotype, 1, 0)
	// Fitness is max (0) when there are 10 'a's
	fitness := func(f *ArrayGenotype[rune]) float64 {
		total := 0.0
		for i := range f.values {
			if f.values[i] == 'a' {
				total += 1
			}
		}
		return -math.Abs(10 - total)
	}
	testWithFitnessFunc(t, fitness, pop)
}

func TestBoolGenotype(t *testing.T) {
	valueGen := NewGeneratorChoices([]bool{false, true})
	mut := NewArrayMutationGeneratorReplace(valueGen, 0.1)
	newGenotype := func() *ArrayGenotype[bool] { return NewArrayGenotype(10, valueGen) }
	pop := setupArrayTestStuff(mut, newGenotype, 2, 1)
	// Fitness is max (0) when there are 10 'true's
	fitness := func(f *ArrayGenotype[bool]) float64 {
		total := 0.0
		for i := range f.values {
			if f.values[i] {
				total += 1
			}
		}
		return -math.Abs(10 - total)
	}
	testWithFitnessFunc(t, fitness, pop)
}
