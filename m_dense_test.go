package goevo

import (
	"testing"
)

func setupDenseTestStuff(numIn, numOut int) Population[*DenseGenotype] {
	selec := NewTournamentSelection[*DenseGenotype](3)

	diffGen := NewGeneratorNormal(0.0, 0.1)
	add := func(old, new float64) float64 { return old + new }
	mut := NewDenseMutationUniform(diffGen, add, 0.1, diffGen, add, 0.1)

	crs := &denseCrossoverUniform{
		parents: 2,
	}

	reprod := NewTwoPhaseReproduction(crs, mut)

	gen := NewGeneratorNormal(0.0, 0.5)

	var pop Population[*DenseGenotype] = NewSimplePopulation(func() *DenseGenotype {
		return NewDenseGenotype([]int{numIn, 5, numOut}, Linear, Relu, Sigmoid, gen, gen)
	}, 100, selec, reprod)

	return pop
}

func TestDenseXOR(t *testing.T) {
	testWithXORDataset(t, setupDenseTestStuff(3, 1), nil)
}
