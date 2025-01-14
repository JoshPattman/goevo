package goevo

import (
	"fmt"
	"testing"
)

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
