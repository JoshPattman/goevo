package goevo

// This file contains the implementation of a hyperneat phenotype that is flattened into a standard dense neural net.
// You can create a LayeredHyperPhenotype from either a LayeredSubstrate1D or a LayeredSubstrateND.
// I plan to add more implementations of different HyperNEAT phenotypes in the future.

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// LayeredSubstrate1D stores information about a substrate for a LayeredHyperPhenotype.
// It follows the structure of a dense neural network, where there are multiple layers, each layer having an activation and a number of nodes.
type LayeredSubstrate1D struct {
	Dimensions  []int        `json:"dimensions"`
	Activations []Activation `json:"activations"`
}

// LayeredHyperPhenotype is a HyperNEAT phenotype created with a substrate composed of a number of n dimensional layers.
type LayeredHyperPhenotype struct {
	weights     []*mat.Dense // For efficiency, we can flatten the n dimensional layers into flat layers and the code works exactly the same.
	activations []func(float64) float64
}

// NewLayeredHyperPhenotype1D creates a new LayeredHyperPhenotype from a LayeredSubstrate1D and a CPPN.
// The CPPN should have 5 inputs, the first two being the layer and node of the input, the second two being the layer and node of the output, and the last being a bias node.
// The CPPN should have 1 output, the weight of the connection.
func NewLayeredHyperPhenotype1D(substrate *LayeredSubstrate1D, cppn *Phenotype) *LayeredHyperPhenotype {
	weights := make([]*mat.Dense, len(substrate.Dimensions)-1)
	acfuncs := make([]func(float64) float64, len(substrate.Dimensions))
	for i, a := range substrate.activations {
		acfuncs[i] = activationMap[a]
	}
	for layer := 0; layer < len(substrate.Dimensions)-1; layer++ {
		weights[layer] = mat.NewDense(substrate.Dimensions[layer+1], substrate.Dimensions[layer]+1, nil)
		for out := 0; out < substrate.Dimensions[layer+1]; out++ {
			for inp := 0; inp < substrate.Dimensions[layer]+1; inp++ {
				layerDivisor := float64(len(substrate.Dimensions) - 1)
				inpDivisor, outDivisor := float64(substrate.Dimensions[layer]+1), float64(substrate.Dimensions[layer+1])
				weight := cppn.Forward([]float64{
					float64(layer) / layerDivisor, float64(inp) / inpDivisor,
					float64(layer+1) / layerDivisor, float64(out) / outDivisor,
					1,
				})[0]
				if math.IsNaN(weight) {
					panic("cppn generated a nan value")
				}
				weights[layer].Set(out, inp, weight)
			}
		}
	}
	return &LayeredHyperPhenotype{
		weights:     weights,
		activations: acfuncs,
	}
}

// Forward performs a forward pass on the LayeredHyperPhenotype, returning the output of the network.
func (p *LayeredHyperPhenotype) Forward(inputs []float64) []float64 {
	buf := mat.NewVecDense(len(inputs)+1, append(inputs, 1))
	for i := 0; i < buf.Len(); i++ {
		buf.SetVec(i, p.activations[0](buf.AtVec(i)))
	}
	for l := range p.weights {
		eouts, _ := p.weights[l].Dims()
		temp := mat.NewVecDense(eouts, nil)
		temp.MulVec(p.weights[l], buf)
		for i := 0; i < eouts; i++ {
			temp.SetVec(i, p.activations[l+1](temp.AtVec(i)))
		}
		if l == len(p.weights)-1 {
			// Dont add a bias node on the last layer
			buf = temp
		} else {
			buf = mat.NewVecDense(eouts+1, append(temp.RawVector().Data, 1))
		}
	}
	return buf.RawVector().Data
}
