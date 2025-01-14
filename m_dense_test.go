package goevo

import (
	"fmt"
	"testing"
)

func setupDenseTestStuff(numIn, numOut int) Population[*DenseGenotype] {
	selec := &TournamentSelection[*DenseGenotype]{
		TournamentSize: 3,
	}

	mut := &DenseMutationStd{
		WeightStd:    0.1,
		WeightChance: 0.1,
		WeightMax:    10,
		BiasStd:      0.1,
		BiasChance:   0.1,
		BiasMax:      10,
	}

	crs := &DenseCrossoverUniform{
		Parents: 2,
	}

	reprod := NewTwoPhaseReproduction(crs, mut)

	gen := &NormalGenerator[float64]{
		Std: 0.5,
	}

	var pop Population[*DenseGenotype] = NewSimplePopulation(func() *DenseGenotype {
		return NewDenseGenotype([]int{numIn, 5, numOut}, Linear, Relu, Sigmoid, gen, gen)
	}, 100, selec, reprod)

	return pop
}

func TestNewDenseGenotype(t *testing.T) {
	wg := &NormalGenerator[float64]{
		Std: 0.5,
	}
	bg := &NormalGenerator[float64]{
		Std: 0.1,
	}
	g1 := NewDenseGenotype([]int{5, 4, 3}, Linear, Relu, Softmax, wg, bg)
	fmt.Println(g1.Forward([]float64{0, 1, 2, 3, 4}))

	m := &DenseMutationStd{
		WeightStd:    0.1,
		WeightChance: 1,
		WeightMax:    2,
		BiasStd:      0.1,
		BiasChance:   1,
		BiasMax:      1,
	}
	g2 := Clone(g1)
	m.Mutate(g2)
	fmt.Println(g2.Forward([]float64{0, 1, 2, 3, 4}))

	c := &DenseCrossoverUniform{
		Parents: 2,
	}
	g3 := c.Crossover([]*DenseGenotype{g1, g2})
	fmt.Println(g3.Forward([]float64{0, 1, 2, 3, 4}))
}

func TestDenseXOR(t *testing.T) {
	testWithXORDataset(t, setupDenseTestStuff(3, 1), nil)
}
