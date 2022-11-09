package goevo

type PhenotypeConnection struct {
	To     int
	Weight float64
}
type Phenotype struct {
	memory      []float64
	activations [](func(float64) float64)
	conns       [][]PhenotypeConnection
	numIn       int
	numOut      int
}

func NewPhenotype(g *Genotype) *Phenotype {
	mem := make([]float64, len(g.Neurons))
	acts := make([](func(float64) float64), len(g.Neurons))
	conns := make([][]PhenotypeConnection, len(g.Neurons))
	for o, oid := range g.NeuronOrder {
		connectedNeurons := make([]PhenotypeConnection, 0)
		for _, s := range g.Synapses {
			if s.From == oid {
				toOrder, _ := g.GetNeuronOrder(s.To)
				connectedNeurons = append(connectedNeurons, PhenotypeConnection{toOrder, s.Weight})
			}
		}
		conns[o] = connectedNeurons
		acts[o] = activationMap[g.Neurons[oid].Activation]
	}
	return &Phenotype{
		memory:      mem,
		activations: acts,
		conns:       conns,
		numIn:       g.NumIn,
		numOut:      g.NumOut,
	}
}

func (p *Phenotype) Forward(inputs []float64) []float64 {
	if len(inputs) != p.numIn {
		panic("not correct number of inputs")
	}
	for i := range p.memory {
		if i < p.numIn {
			p.memory[i] = inputs[i]
		} else {
			p.memory[i] = 0
		}
	}
	for ni := range p.memory {
		p.memory[ni] = p.activations[ni](p.memory[ni])
		for _, c := range p.conns[ni] {
			p.memory[c.To] += c.Weight * p.memory[ni]
		}
	}
	output := make([]float64, p.numOut)
	copy(output, p.memory[len(p.memory)-p.numOut:])
	return output
}
