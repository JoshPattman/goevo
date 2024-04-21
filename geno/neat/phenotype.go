package neat

import "github.com/JoshPattman/goevo"

type phenotypeConnection struct {
	toIdx int
	w     float64
}

var _ goevo.Forwarder = &Phenotype{}

// Phenotype is a phenotype for a NEAT genotype.
// It conceptually represents a neural network, built according to the instructions in the NEATGenotype (DNA).
// Once built, the Phenotype can be used to forward propagate inputs through the network,
// but it cannot be modified though mutation or corss-over.
type Phenotype struct {
	numIn            int
	numOut           int
	accumulators     []float64
	lastAccumulators []float64
	activations      []goevo.Activation
	forwardWeights   [][]phenotypeConnection
	recurrentWeights [][]phenotypeConnection
}

// Build a NEATPhenotype from a NEATGenotype.
func (g *Genotype) Build() goevo.Forwarder {
	accs := make([]float64, len(g.neuronOrder))
	laccs := make([]float64, len(g.neuronOrder))
	acts := make([]goevo.Activation, len(g.neuronOrder))
	fwdWeights := make([][]phenotypeConnection, len(g.neuronOrder))
	recurrentWeights := make([][]phenotypeConnection, len(g.neuronOrder))
	for no, nid := range g.neuronOrder {
		acts[no] = g.activations[nid]
		fwdWeights[no] = make([]phenotypeConnection, 0)
		recurrentWeights[no] = make([]phenotypeConnection, 0)
	}
	for sid, w := range g.weights {
		ep := g.synapseEndpointLookup[sid]
		oa, ob := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
		if ob > oa {
			fwdWeights[oa] = append(fwdWeights[oa], phenotypeConnection{ob, w})
		} else {
			recurrentWeights[oa] = append(recurrentWeights[oa], phenotypeConnection{ob, w})
		}
	}
	return &Phenotype{
		numIn:            g.numInputs,
		numOut:           g.numOutputs,
		accumulators:     accs,
		lastAccumulators: laccs,
		activations:      acts,
		forwardWeights:   fwdWeights,
		recurrentWeights: recurrentWeights,
	}
}

// Forward propagate inputs through the network, returning the resulting outputs.
func (p *Phenotype) Forward(x []float64) []float64 {
	if len(x) != p.numIn {
		panic("incorrect number of inputs")
	}
	// Reset accumulators to default vals (remember what they were before incase we need recurrent connections)
	if len(p.recurrentWeights) > 0 { // For efficiency
		copy(p.lastAccumulators, p.accumulators)
	}
	for i := 0; i < len(p.accumulators); i++ {
		if i < len(x) {
			p.accumulators[i] = x[i]
		} else {
			p.accumulators[i] = 0
		}
	}
	if len(p.recurrentWeights) > 0 { // For efficiency
		// Apply recurrent connections (does not matter what order we do this in)
		for i := 0; i < len(p.accumulators); i++ {
			for _, w := range p.recurrentWeights[i] {
				p.accumulators[w.toIdx] += w.w * p.lastAccumulators[i]
			}
		}
	}
	// Apply forward connections
	for i := 0; i < len(p.accumulators); i++ {
		p.accumulators[i] = goevo.Activate(p.accumulators[i], p.activations[i])
		for _, w := range p.forwardWeights[i] {
			p.accumulators[w.toIdx] += w.w * p.accumulators[i]
		}
	}
	// Copy output array to avoid modification of internal state
	outs := make([]float64, p.numOut)
	copy(outs, p.accumulators[len(p.accumulators)-p.numOut:])
	// Reuturn
	return outs
}

// Reset will clear the recurrent memories of the phenotype.
func (p *Phenotype) Reset() {
	for i := range p.accumulators {
		p.accumulators[i] = 0
	}
}
