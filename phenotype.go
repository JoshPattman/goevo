package goevo

import (
	"fmt"
	"math"
)

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

// Data type representing a phenotype, a sort of compiled genotype
type Phenotype struct {
	memory         []float64
	activations    [](func(float64) float64)
	conns          [][]PhenotypeConnection
	recurrentConns [][]RecurrentPhenotypeConnection
	numIn          int
	numOut         int
}

// Create a phenotype from genotype `g`
func NewPhenotype(g *Genotype) *Phenotype {
	mem := make([]float64, len(g.Neurons))
	acts := make([](func(float64) float64), len(g.Neurons))
	conns := make([][]PhenotypeConnection, len(g.Neurons))
	recurrentConns := make([][]RecurrentPhenotypeConnection, len(g.Neurons))

	for n := range conns {
		conns[n] = make([]PhenotypeConnection, 0)
		recurrentConns[n] = make([]RecurrentPhenotypeConnection, 0)
		acts[n] = activationMap[g.Neurons[g.NeuronOrder[n]].Activation]
	}

	neuronOrderCache := make(map[int]int)

	for _, s := range g.Synapses {
		var fromOrder int
		var toOrder int
		if val, ok := neuronOrderCache[s.From]; ok {
			fromOrder = val
		} else {
			fromOrder, _ = g.GetNeuronOrder(s.From)
			neuronOrderCache[s.From] = fromOrder
		}
		if val, ok := neuronOrderCache[s.To]; ok {
			toOrder = val
		} else {
			toOrder, _ = g.GetNeuronOrder(s.To)
			neuronOrderCache[s.To] = toOrder
		}
		if fromOrder < toOrder {
			conns[fromOrder] = append(conns[fromOrder], PhenotypeConnection{toOrder, s.Weight})
		} else {
			recurrentConns[toOrder] = append(recurrentConns[toOrder], RecurrentPhenotypeConnection{fromOrder, s.Weight})
		}
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

// Do a forward pass with some input data for the phenotype, returning the output of the network. This also will take into account any memory left over from previous calls
func (p *Phenotype) Forward(inputs []float64) []float64 {
	if len(inputs) != p.numIn {
		panic(fmt.Errorf("not correct number of inputs: expected %v but got %v", p.numIn, len(inputs)))
	}
	for _, x := range inputs {
		if math.IsNaN(x) {
			panic("input to phenotype contains nan")
		}
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

// Clear any memory retained from previous calls to `p.Forward`
func (p *Phenotype) ClearRecurrentMemory() {
	for i := range p.memory {
		p.memory[i] = 0
	}
}
