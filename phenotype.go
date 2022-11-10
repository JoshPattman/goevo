package goevo

// Data type for a phenotype connection
type PhenotypeConnection struct {
	To     int
	Weight float64
}

// Data type for a recurrent phenotype connection
type RecurrentPhenotypeConnection struct {
	From   int
	Weight float64
}

// Data type representing a phenotype (a bit like an instance of a genotype)
type Phenotype struct {
	memory         []float64
	activations    [](func(float64) float64)
	conns          [][]PhenotypeConnection
	recurrentConns [][]RecurrentPhenotypeConnection
	numIn          int
	numOut         int
}

// Create a phenotype from a genotype
func NewPhenotype(g *Genotype) *Phenotype {
	mem := make([]float64, len(g.Neurons))
	acts := make([](func(float64) float64), len(g.Neurons))
	conns := make([][]PhenotypeConnection, len(g.Neurons))
	recurrentConns := make([][]RecurrentPhenotypeConnection, len(g.Neurons))
	for o, oid := range g.NeuronOrder {
		connectedNeurons := make([]PhenotypeConnection, 0)
		recurrentConnectedNeurons := make([]RecurrentPhenotypeConnection, 0)
		for _, s := range g.Synapses {
			if s.From == oid {
				fromOrder, _ := g.GetNeuronOrder(s.From)
				toOrder, _ := g.GetNeuronOrder(s.To)
				if fromOrder < toOrder {
					connectedNeurons = append(connectedNeurons, PhenotypeConnection{toOrder, s.Weight})
				}
			} else if s.To == oid {
				fromOrder, _ := g.GetNeuronOrder(s.From)
				toOrder, _ := g.GetNeuronOrder(s.To)
				if fromOrder > toOrder {
					//recurrent
					recurrentConnectedNeurons = append(recurrentConnectedNeurons, RecurrentPhenotypeConnection{fromOrder, s.Weight})
				}
			}
		}
		conns[o] = connectedNeurons
		recurrentConns[o] = recurrentConnectedNeurons
		acts[o] = activationMap[g.Neurons[oid].Activation]
	}
	return &Phenotype{
		memory:         mem,
		activations:    acts,
		conns:          conns,
		recurrentConns: recurrentConns,
		numIn:          g.NumIn,
		numOut:         g.NumOut,
	}
}

// Do a forward pass with some input data for the phenotype, returning the output of the network
func (p *Phenotype) Forward(inputs []float64) []float64 {
	if len(inputs) != p.numIn {
		panic("not correct number of inputs")
	}
	for i := range p.memory {
		if i < p.numIn {
			p.memory[i] += inputs[i]
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
	for ni := range p.memory {
		p.memory[ni] = 0
		for _, c := range p.recurrentConns[ni] {
			p.memory[ni] += p.memory[c.From] * c.Weight
		}
	}
	return output
}

func (p *Phenotype) ClearRecurrentMemory() {
	for i := range p.memory {
		p.memory[i] = 0
	}
}
