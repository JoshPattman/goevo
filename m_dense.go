package goevo

import "gonum.org/v1/gonum/mat"

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
	gn := &DenseGenotype{
		weights:          make([]*mat.Dense, len(d.weights)),
		biases:           make([]*mat.VecDense, len(d.biases)),
		buffers:          make([]*mat.VecDense, len(d.buffers)),
		inputActivation:  d.inputActivation,
		hiddenActivation: d.hiddenActivation,
		outputActivation: d.outputActivation,
	}

	for bi := range d.buffers {
		gn.biases[bi] = mat.VecDenseCopyOf(d.biases[bi])
		gn.buffers[bi] = mat.VecDenseCopyOf(d.buffers[bi])
	}
	for wi := range d.weights {
		gn.weights[wi] = mat.DenseCopyOf(d.weights[wi])
	}
	return gn
}

// denseMutationUniform is a type of mutation for dense genotypes.
// For each weight and bias, with a certain chance, it mutates the value within a normal distribution.
// It then caps all values at the maximum value.
type denseMutationUniform struct {
	genWeights     Generator[float64]
	combineWeights func(old, new float64) float64
	genBiases      Generator[float64]
	combineBiases  func(old, new float64) float64
	weightsChance  float64
	biasesChance   float64
}

func NewDenseMutationUniform(
	genWeights Generator[float64],
	combineWeights func(old, new float64) float64,
	weightsChance float64,
	genBiases Generator[float64],
	combineBiases func(old, new float64) float64,
	biasesChance float64,
) Mutation[*DenseGenotype] {
	if genWeights == nil || genBiases == nil {
		panic("cannot have nil generator")
	}
	if combineWeights == nil || combineBiases == nil {
		panic("cannot have nil combine func")
	}
	if weightsChance < 0 || weightsChance > 1 || biasesChance < 0 || biasesChance > 1 {
		panic("cannot have chance out of range 0-1")
	}
	return &denseMutationUniform{
		genWeights:     genWeights,
		genBiases:      genBiases,
		combineWeights: combineWeights,
		combineBiases:  combineBiases,
		weightsChance:  weightsChance,
		biasesChance:   biasesChance,
	}
}

// Mutate implements Mutation.
func (m *denseMutationUniform) Mutate(g *DenseGenotype) {
	for _, w := range g.weights {
		rs, cs := w.Dims()
		for r := range rs {
			for c := range cs {
				v := w.At(r, c)
				v = m.combineWeights(v, m.genWeights.Next())
				w.Set(r, c, v)
			}
		}
	}
	for _, b := range g.biases {
		rs := b.Len()
		for r := range rs {
			v := b.AtVec(r)
			v = m.combineBiases(v, m.genBiases.Next())
			b.SetVec(r, v)
		}
	}
}

// denseCrossoverUniform is a type of crossover for dense genotypes.
// For each weight and bias, it chooses randomly from one of its parents.
// The number of parents is a parameter.
type denseCrossoverUniform struct {
	parents int
}

func NewDenseCrossoverUniform(parents int) Crossover[*DenseGenotype] {
	if parents <= 0 {
		panic("must have at least one parent")
	}
	return &denseCrossoverUniform{
		parents: parents,
	}
}

// Crossover implements Crossover.
func (c *denseCrossoverUniform) Crossover(parents []*DenseGenotype) *DenseGenotype {
	if len(parents) != c.parents {
		panic("incorrect number of parents")
	}
	if c.parents <= 0 {
		panic("must have at least one parent")
	}
	g := Clone(parents[0])
	for _, p := range parents {
		if len(p.weights) != len(g.weights) {
			panic("inconsistent parent num layers")
		}
	}
	for wi, w := range g.weights {
		wr, wc := w.Dims()
		pws := make([]mutMat, len(parents))
		for pi, p := range parents {
			pr, pc := p.weights[wi].Dims()
			if wr != pr || wc != pc {
				panic("incorrect weight sizes")
			}
			pws[pi] = p.weights[wi]
		}
		randomChoiceMatrix(w, pws)
	}
	for bi, b := range g.biases {
		br := b.Len()
		pbs := make([]mutMat, len(parents))
		for pi, p := range parents {
			pr := p.biases[bi].Len()
			if br != pr {
				panic("incorrect bias sizes")
			}
			pbs[pi] = &mutVecWrapper{p.biases[bi]}
		}
		randomChoiceMatrix(&mutVecWrapper{b}, pbs)
	}
	return g
}

// NumParents implements Crossover.
func (c *denseCrossoverUniform) NumParents() int {
	return c.parents
}
