package goevo

import (
	"gonum.org/v1/gonum/mat"
)

// NewLayeredSubstrate creates a new LayeredSubstrate with the given parameters.
// layeredNeuronPositions is a slice of slices of positions, where each slice of positions represents a layer.
// biasNeuronPosition is the position of the bias neuron, which is placed in the previous layer. It should have same ndims as each position in the layerNeuronPositions.
// It will panic on invalid input.
func NewLayeredSubstrate(layeredNeuronPositions [][]Pos, layerActivations []Activation, biasNeuronPosition Pos) *LayeredSubstrate {
	if len(layeredNeuronPositions) < 2 {
		panic("LayeredSubstrate must have at least an input and ouput layer")
	}
	for l := range layeredNeuronPositions {
		if len(layeredNeuronPositions[l]) == 0 {
			panic("all layers must have at least one neuron")
		}
	}
	dims := len(layeredNeuronPositions[0][0])
	for l := range layeredNeuronPositions {
		for p := range layeredNeuronPositions[l] {
			if len(layeredNeuronPositions[l][p]) != dims {
				panic("all posses must have the same ndims")
			}
		}
	}
	if len(biasNeuronPosition) != dims {
		panic("bias neuraon must have same ndims as other neurons")
	}
	if len(layeredNeuronPositions) != len(layerActivations) {
		panic("there must be exactly one activation for each layer")
	}

	return &LayeredSubstrate{
		LayeredNeuronPositions: layeredNeuronPositions,
		BiasNeuronPosition:     biasNeuronPosition,
		LayerActivations:       layerActivations,
	}
}

// LayeredSubstrate stores information about a substrate for a LayeredHyperPhenotype.
// It follows the structure of a dense neural network, where there are multiple layers, each layer having an activation and a number of nodes.
// However, when generation of weights using the CPPN, the nodes are not arranged in a array, but rather are scattered throughout the substrate in user specified positions.
// When adding the bias, the bias neuron is always placed at the position of BiasNeuronPosition but in the previous layer.
type LayeredSubstrate struct {
	LayeredNeuronPositions [][]Pos      `json:"layered_neuron_positions"`
	BiasNeuronPosition     Pos          `json:"bias_neuron_position"`
	LayerActivations       []Activation `json:"layer_activations"`
}

// CPNNInputsOutputs returns the number of inputs and outputs the CPPN should have for this substrate.
func (s *LayeredSubstrate) CPNNInputsOutputs() (int, int) {
	return (s.Dimensions()+1)*2 + 1, 1
}

// NewPhenotype creates a new LayeredHyperPhenotype from the CPPN using this substrate.
func (s *LayeredSubstrate) NewPhenotype(cppn Forwarder) *LayeredHyperPhenotype {
	numLayers := len(s.LayeredNeuronPositions)
	weights := make([]*mat.Dense, numLayers-1)
	activations := make([]func(float64) float64, numLayers)

	for layer := range activations {
		activations[layer] = activationMap[s.LayerActivations[layer]]
	}

	for srcLayer := 0; srcLayer < numLayers-1; srcLayer++ {
		tarLayer := srcLayer + 1
		srcNum, tarNum := len(s.LayeredNeuronPositions[srcLayer])+1, len(s.LayeredNeuronPositions[tarLayer]) // Add one to srcNum for bias
		weights[srcLayer] = mat.NewDense(tarNum, srcNum, nil)
		for src := 0; src < srcNum; src++ {
			var srcPos Pos
			if src == srcNum-1 {
				srcPos = s.BiasNeuronPosition
			} else {
				srcPos = s.LayeredNeuronPositions[srcLayer][src]
			}
			for tar := 0; tar < tarNum; tar++ {
				tarPos := s.LayeredNeuronPositions[tarLayer][tar]
				srcPosWithLayer, tarPosWithLayer := append(srcPos, float64(srcLayer)), append(tarPos, float64(tarLayer))
				cppnInputs := make([]float64, (s.Dimensions()+1)*2+1)
				for i := 0; i < s.Dimensions()+1; i++ { // add one to dimensions as dimensions does not include layer
					cppnInputs[i] = srcPosWithLayer[i]
					cppnInputs[i+s.Dimensions()+1] = tarPosWithLayer[i]
				}
				cppnInputs[(s.Dimensions()+1)*2] = 1 // Bias
				// The structure of cppnInputs is (srcPos.X, srcPos.Y, tarPos.X, tarPos.Y, bias) but can be more or less than just X and Y
				// The output of the CPPN is the weight of the synapse
				weight := cppn.Forward(cppnInputs)[0]
				weights[srcLayer].Set(tar, src, weight)
			}
		}
	}
	return &LayeredHyperPhenotype{
		weights:     weights,
		activations: activations,
	}
}

// Dimensions returns the number of dimensions each layer has. If each layer has one axis (P(0.5)), this will return 1
func (s *LayeredSubstrate) Dimensions() int {
	return len(s.LayeredNeuronPositions[0][0])
}

// LayeredHyperPhenotype is a HyperNEAT phenotype created with a substrate composed of a number of n dimensional layers.
type LayeredHyperPhenotype struct {
	weights     []*mat.Dense // For efficiency, we can flatten the n dimensional layers into flat layers and the code works exactly the same.
	activations []func(float64) float64
}

// Forward performs a forward pass on the LayeredHyperPhenotype, returning the output of the network.
func (p *LayeredHyperPhenotype) Forward(inputs []float64) []float64 {
	buf := mat.NewVecDense(len(inputs)+1, append(inputs, 1))
	for i := 0; i < buf.Len()-1; i++ { // Dont want to activate the bias node
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
