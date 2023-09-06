package goevo

import (
	"gonum.org/v1/gonum/mat"
)

type HyperPhenotype struct {
	Weights     []*mat.Dense
	Activations []func(float64) float64
}

// Creates a hyperneat phenotype which has the same structure as a standard dene neural net
// The cppn should have 5 inputs: (la, ia, lb, ib, bias)
func NewFlatHyperPhenotype(dimensions []int, activations []Activation, cppn *Phenotype) *HyperPhenotype {
	weights := make([]*mat.Dense, len(dimensions)-1)
	acfuncs := make([]func(float64) float64, len(dimensions))
	for i, a := range activations {
		acfuncs[i] = activationMap[a]
	}
	for layer := 0; layer < len(dimensions)-1; layer++ {
		weights[layer] = mat.NewDense(dimensions[layer+1], dimensions[layer]+1, nil)
		for out := 0; out < dimensions[layer+1]; out++ {
			for inp := 0; inp < dimensions[layer]+1; inp++ {
				layerDivisor := float64(len(dimensions) - 1)
				inpDivisor, outDivisor := float64(dimensions[layer]), float64(dimensions[layer+1]-1) // -1 for both of these to mean that the first node is 0, and the last is 1
				weight := cppn.Forward([]float64{
					float64(layer) / layerDivisor, float64(inp) / inpDivisor,
					float64(layer+1) / layerDivisor, float64(out) / outDivisor,
					1,
				})[0]
				weights[layer].Set(out, inp, weight)
			}
		}
	}
	return &HyperPhenotype{
		Weights:     weights,
		Activations: acfuncs,
	}
}

func (p *HyperPhenotype) Forward(inputs []float64) []float64 {
	buf := mat.NewVecDense(len(inputs)+1, append(inputs, 1))
	for l := range p.Weights {
		eouts, _ := p.Weights[l].Dims()
		temp := mat.NewVecDense(eouts, nil)
		temp.MulVec(p.Weights[l], buf)
		for i := 0; i < eouts; i++ {
			temp.SetVec(i, p.Activations[l](temp.AtVec(i)))
		}
		if l == len(p.Weights)-1 {
			// Dont add a bias node on the last layer
			buf = temp
		} else {
			buf = mat.NewVecDense(eouts+1, append(temp.RawVector().Data, 1))
		}
	}
	return buf.RawVector().Data
}
