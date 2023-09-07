package goevo

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// LayeredSubstrateScattered stores information about a substrate for a LayeredHyperPhenotype.
// It follows the structure of a dense neural network, where there are multiple layers, each layer having an activation and a number of nodes.
// However, when using the CPPN, the nodes are not arranged in a array, but rather are scattered throughout the substrate in user specified positions.
type LayeredSubstrate struct {
	Neurons     [][]Pos      `json:"neurons"`
	Activations []Activation `json:"activations"`
}

// LayeredHyperPhenotype is a HyperNEAT phenotype created with a substrate composed of a number of n dimensional layers.
type LayeredHyperPhenotype struct {
	weights     []*mat.Dense // For efficiency, we can flatten the n dimensional layers into flat layers and the code works exactly the same.
	activations []func(float64) float64
}

// NewLayeredHyperPhenotype creates a new LayeredHyperPhenotype from a LayeredSubstrate and a CPPN.
// The CPPN should have 5 inputs, the first two being the layer and node of the input, the second two being the layer and node of the output, and the last being a bias node.
// The CPPN should have 1 output, the weight of the connection.
func NewLayeredHyperPhenotype(substrate *LayeredSubstrate, cppn *Phenotype) *LayeredHyperPhenotype {
	weights := make([]*mat.Dense, len(substrate.Neurons)-1)
	acfuncs := make([]func(float64) float64, len(substrate.Neurons))
	for i, a := range substrate.Activations {
		acfuncs[i] = activationMap[a]
	}
	for layer := 0; layer < len(substrate.Neurons)-1; layer++ {
		weights[layer] = mat.NewDense(len(substrate.Neurons[layer+1]), len(substrate.Neurons[layer])+1, nil)
		for out := 0; out < len(substrate.Neurons[layer+1]); out++ {
			for inp := 0; inp < len(substrate.Neurons[layer])+1; inp++ {
				layerDivisor := float64(len(substrate.Neurons) - 1)
				cppnInputs := []float64{float64(layer) / layerDivisor, float64(layer+1) / layerDivisor}
				cppnInputs = append(cppnInputs, substrate.Neurons[layer][inp]...)
				cppnInputs = append(cppnInputs, substrate.Neurons[layer+1][out]...)
				cppnInputs = append(cppnInputs, 1)
				weight := cppn.Forward(cppnInputs)[0]
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
