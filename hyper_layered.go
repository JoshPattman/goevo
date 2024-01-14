package goevo

import (
	"encoding/json"
	"io"
	"math"
	"os"

	"gonum.org/v1/gonum/mat"
)

// LayeredSubstrate stores information about a substrate for a LayeredHyperPhenotype.
// It follows the structure of a dense neural network, where there are multiple layers, each layer having an activation and a number of nodes.
// However, when generation of weights using the CPPN, the nodes are not arranged in a array, but rather are scattered throughout the substrate in user specified positions.
// When adding the bias, the bias neuron is always placed at the position of BiasNeuronPosition but in the previous layer.
type LayeredSubstrate struct {
	NodeLateralPositions [][]Pos      `json:"node_lateral_positions"`
	LayerPositions       []Pos        `json:"layer_positions"`
	BiasLateralPosition  Pos          `json:"bias_lateral_position"`
	LayerActivations     []Activation `json:"layer_activations"`
}

func NewLayeredSubstrateEmpty() *LayeredSubstrate {
	return &LayeredSubstrate{}
}

// NewLayeredSubstrate creates a new substrate for creating a layered hyper phenotype, using a the manually provided node positions.
//
//   - nodeLateralPositions: The lateral positions of each neuron in each layer
//   - layerPositions: The position of each layer, will be appended onto the end of each node pos
//   - layerActivations: An activation for each layer
//   - biasLateralPosition: The lateral position of the bias neuron, bias neuron for each layer will always be placed in th eprevious layer
func NewLayeredSubstrate(nodeLateralPositions [][]Pos, layerPositions []Pos, layerActivations []Activation, biasLateralPosition Pos) *LayeredSubstrate {
	if len(nodeLateralPositions) < 2 {
		panic("LayeredSubstrate must have at least an input and ouput layer")
	}
	if len(nodeLateralPositions) != len(layerPositions) || len(nodeLateralPositions) != len(layerActivations) {
		panic("incorrect lengths. p.s. write better error message")
	}
	for l := range nodeLateralPositions {
		if len(nodeLateralPositions[l]) == 0 {
			panic("all layers must have at least one neuron")
		}
	}
	dims := len(nodeLateralPositions[0][0]) + len(layerPositions[0])
	for l := range nodeLateralPositions {
		for p := range nodeLateralPositions[l] {
			if len(nodeLateralPositions[l][p])+len(layerPositions[l]) != dims {
				panic("all posses must have the same ndims")
			}
		}
	}
	if len(biasLateralPosition)+len(layerPositions[0]) != dims {
		panic("bias neuron must have same ndims as other neurons")
	}

	return &LayeredSubstrate{
		NodeLateralPositions: nodeLateralPositions,
		LayerPositions:       layerPositions,
		BiasLateralPosition:  biasLateralPosition,
		LayerActivations:     layerActivations,
	}
}

// NewProceduralSinLayeredSubstrate creates a new substrate for creating a layered hyper phenotype, using generated positions for nodes.
// Instead of the programmer manually placing nodes in ndimensional space, this function procedurally generates the positions of the nodes, using a sin-based positional encoding.
//
//   - layerNeuronCounts: The topology of the network. E.g. [3, 5, 4, 1] would give 3 input, 5 hidden1, 4 hidden2, and 1 output
//   - layerActivations: An activation for each layer
//   - nodeEncodingDim: The dimension of the lateral position of each node
//   - layerEncodingDim: The dinension of the layer position
//   - nodeEncodingMaxFreq: The max frequency of sin waves for the node encoding, 4 seems to work well
//   - layerEncodingMaxFreq: The max frequency of sin waves for the layer encoding, 4 seems to work well
func NewProceduralSinLayeredSubstrate(layerNeuronCounts []int, layerActivations []Activation, nodeEncodingDim, layerEncodingDim int, nodeEncodingMaxFreq, layerEncodingMaxFreq float64) *LayeredSubstrate {
	lateralPosses := make([][]Pos, len(layerNeuronCounts))
	layerPosses := make([]Pos, len(layerActivations))
	for l := range lateralPosses {
		layerPosses[l] = positionalEncoding(float64(l)/float64(len(layerNeuronCounts)), layerEncodingMaxFreq, layerEncodingDim, 0.5) // 0.5 is a good offset for all uses i think
		lateralPosses[l] = make([]Pos, layerNeuronCounts[l])
		for n := range lateralPosses[l] {
			lateralPosses[l][n] = positionalEncoding(float64(n)/float64(layerNeuronCounts[l]), nodeEncodingMaxFreq, nodeEncodingDim, 0.5)
		}
	}
	return NewLayeredSubstrate(lateralPosses, layerPosses, layerActivations, make(Pos, nodeEncodingDim))
}

func positionalEncoding(p float64, maximumMult float64, dims int, offset float64) []float64 {
	mult := math.Pow(maximumMult, 1/float64(dims-1))
	totalOffset := 0.0
	v := 1.0
	encoding := make([]float64, dims)
	for i := 0; i < dims; i++ {
		encoding[i] = math.Sin(v*p*2*math.Pi + totalOffset)
		v *= mult
		totalOffset += offset
	}
	return encoding
}

// CPNNInputsOutputs returns the number of inputs and outputs the CPPN should have for this substrate.
func (s *LayeredSubstrate) CPNNInputsOutputs() (int, int) {
	return (s.Dimensions()+1)*2 + 1, 1
}

// BuildPhenotype creates a new LayeredHyperPhenotype from the CPPN using this substrate.
func (s *LayeredSubstrate) BuildPhenotype(cppn Forwarder) *LayeredHyperPhenotype {
	numLayers := len(s.NodeLateralPositions)
	weights := make([]*mat.Dense, numLayers-1)
	activations := make([]func(float64) float64, numLayers)

	for layer := range activations {
		activations[layer] = activationMap[s.LayerActivations[layer]]
	}

	for srcLayer := 0; srcLayer < numLayers-1; srcLayer++ {
		tarLayer := srcLayer + 1
		srcNum, tarNum := len(s.NodeLateralPositions[srcLayer])+1, len(s.NodeLateralPositions[tarLayer]) // Add one to srcNum for bias
		weights[srcLayer] = mat.NewDense(tarNum, srcNum, nil)
		for src := 0; src < srcNum; src++ {
			var srcPos Pos
			if src == srcNum-1 {
				srcPos = append(append(Pos{}, s.BiasLateralPosition...), s.LayerPositions[srcLayer]...)
			} else {
				srcPos = append(append(Pos{}, s.NodeLateralPositions[srcLayer][src]...), s.LayerPositions[srcLayer]...)
			}
			for tar := 0; tar < tarNum; tar++ {
				tarPos := append(append(Pos{}, s.NodeLateralPositions[tarLayer][tar]...), s.LayerPositions[tarLayer]...)
				cppnInputs := make([]float64, (s.Dimensions()+1)*2+1)
				for i := 0; i < s.Dimensions()+1; i++ { // add one to dimensions as dimensions does not include layer
					cppnInputs[i] = srcPos[i]
					cppnInputs[i+s.Dimensions()+1] = tarPos[i]
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
	return len(s.NodeLateralPositions[0][0])
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

func (s *LayeredSubstrate) WriteJson(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	return enc.Encode(s)
}

// WriteJsonFile writes a JSON representation of this Genotype to a file.
func (s *LayeredSubstrate) WriteJsonFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return s.WriteJson(f)
}

func (s *LayeredSubstrate) ReadJson(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(s)
}

func (s *LayeredSubstrate) ReadJsonFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	return s.ReadJson(f)
}
