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
	g := NewDenseGenotype([]int{5, 4, 3}, Linear, Relu, Softmax, wg, bg)

	fmt.Println(g.Forward([]float64{0, 1, 2, 3, 4}))
}
