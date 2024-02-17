package goevo

type Forwarder interface {
	Forward([]float64) []float64
}

type phenotypeConnection struct {
	toIdx int
	w     float64
}

type Phenotype struct {
	numIn            int
	numOut           int
	accumulators     []float64
	activations      []Activation
	forwardWeights   [][]phenotypeConnection
	recurrentWeights [][]phenotypeConnection
}

func (g *Genotype) Build() *Phenotype {
	accs := make([]float64, len(g.neuronOrder))
	acts := make([]Activation, len(g.neuronOrder))
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
		activations:      acts,
		forwardWeights:   fwdWeights,
		recurrentWeights: recurrentWeights,
	}
}

func (p *Phenotype) Forward(x []float64) []float64 {
	if len(x) != p.numIn {
		panic("incorrect number of inputs")
	}
	// Reset accumulators to default vals
	for i := 0; i < len(p.accumulators); i++ {
		if i < len(x) {
			p.accumulators[i] = x[i]
		} else {
			p.accumulators[i] = 0
		}
	}
	for i := 0; i < len(p.accumulators); i++ {
		p.accumulators[i] = activate(p.accumulators[i], p.activations[i])
		for _, w := range p.forwardWeights[i] {
			p.accumulators[w.toIdx] += w.w * p.accumulators[i]
		}
	}
	outs := make([]float64, p.numOut)
	copy(outs, p.accumulators[len(p.accumulators)-p.numOut:])
	return outs
}
