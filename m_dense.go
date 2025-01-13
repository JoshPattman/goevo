package goevo

import "gonum.org/v1/gonum/mat"

var _ Cloneable = &DenseGenotype{}
var _ Forwarder = &DenseGenotype{}

// DenseGenotype is a type of genotype/phenotype that is a dense feed-forward neural network.
type DenseGenotype struct {
	weights          []*mat.Dense
	biases           []*mat.VecDense
	buffers          []*mat.VecDense
	inputActivation  Activation
	hiddenActivation Activation
	outputActivation Activation
}

// NewDenseGenotype creates a new dense genotype of a shape with activations,
// pulling weights from the generator
func NewDenseGenotype(shape []int, input, hidden, output Activation, weights, biases Generator[float64]) *DenseGenotype {
	if len(shape) < 2 {
		panic("cannot have fewer than two layers")
	}
	numWeights := len(shape) - 1
	g := &DenseGenotype{
		weights:          make([]*mat.Dense, numWeights),
		biases:           make([]*mat.VecDense, numWeights+1),
		buffers:          make([]*mat.VecDense, numWeights+1),
		inputActivation:  input,
		hiddenActivation: hidden,
		outputActivation: output,
	}
	for wi := range numWeights {
		r, c := shape[wi+1], shape[wi]
		data := make([]float64, r*c)
		for i := range data {
			data[i] = weights.Next()
		}
		g.weights[wi] = mat.NewDense(r, c, data)
	}
	for bi := range numWeights + 1 {
		r := shape[bi]
		data := make([]float64, r)
		for i := range data {
			data[i] = biases.Next()
		}
		g.biases[bi] = mat.NewVecDense(r, data)
		g.buffers[bi] = mat.NewVecDense(r, nil)
	}
	return g
}

// Forward implements Forwarder.
func (d *DenseGenotype) Forward(input []float64) []float64 {
	if len(input) != d.buffers[0].Len() {
		panic("incorrect input length")
	}
	// Copy input in
	copy(d.buffers[0].RawVector().Data, input)
	// Perform forward pass
	for bi := range len(d.buffers) {
		// Add bias
		d.buffers[bi].AddVec(d.buffers[bi], d.biases[bi])
		// Activate
		var ac Activation
		if bi == 0 {
			ac = d.inputActivation
		} else if bi == len(d.buffers)-1 {
			ac = d.outputActivation
		} else {
			ac = d.hiddenActivation
		}
		ActivateVector(d.buffers[bi], ac)
		// Multiply by weights if not last
		if bi < len(d.buffers)-1 {
			d.buffers[bi+1].MulVec(d.weights[bi], d.buffers[bi])
		}
	}
	// Copy result out
	lastBuf := d.buffers[len(d.buffers)-1]
	result := make([]float64, lastBuf.Len())
	copy(result, lastBuf.RawVector().Data)
	return result
}

// Clone implements Cloneable.
func (d *DenseGenotype) Clone() any {
	panic("unimplemented")
}
