package goevo

import (
	"errors"
	"math/rand"
)

// A type denoting the type of a neuron (input, hidden, output). Must be one of the consts that start with 'Neuron...'
type NeuronType string

const (
	// Input neuron
	NeuronInput NeuronType = "input"
	// Hidden neuron
	NeuronHidden NeuronType = "hidden"
	// Output neuron
	NeuronOutput NeuronType = "output"
)

// A type represeting a genotype synapse
type Synapse struct {
	// The id of the neuron this comes from
	From int `json:"from_id"`
	// The id of the neuron this goes to
	To int `json:"to_id"`
	// The weight of this synapse
	Weight float64 `json:"weight"`
}

// A type representing a genotype neuron
type Neuron struct {
	// The type of this neuron
	Type NeuronType `json:"type"`
	// The activation of this neuron
	Activation Activation `json:"activation"`
}

// A type representing a genotype (effectively DNA) used in the NEAT algorithm
type Genotype struct {
	// The number of input neurons. READONLY
	NumIn int `json:"num_in"`
	// The number of output neurons. READONLY
	NumOut int `json:"num_out"`
	// A record of all neurons stored according to their id. READONLY
	Neurons map[int]*Neuron `json:"neurons"`
	// A record of all synapses stored according to their id. READONLY
	Synapses map[int]*Synapse `json:"synapses"`
	// A record of all neurons ids stored in execution order. READONLY
	NeuronOrder []int `json:"neuron_order"`
	// A record of all neurons orders stored according to neuron id. READONLY
	InverseNeuronOrder map[int]int `json:"inverse_neuron_order"`
}

func NewGenotypeEmpty() *Genotype {
	return &Genotype{
		NumIn:              -1,
		NumOut:             -1,
		Neurons:            make(map[int]*Neuron),
		Synapses:           make(map[int]*Synapse),
		NeuronOrder:        make([]int, 0),
		InverseNeuronOrder: make(map[int]int),
	}
}

// Create a new Genotype. Create all new innovation IDs with the provided Counter `counter`.
// The genotype will have `numIn` input nodes and `numOut` output nodes.
// The input nodes will have activation `inActivation` and the output nodes will have activation `outActivation`
func NewGenotype(counter Counter, numIn, numOut int, inActivation, outActivation Activation) *Genotype {
	neurons := make(map[int]*Neuron)
	neuronOrder := make([]int, numIn+numOut)
	inverseNeuronOrder := make(map[int]int)
	for i := 0; i < numIn; i++ {
		id := counter.Next()
		neurons[id] = &Neuron{
			Type:       NeuronInput,
			Activation: inActivation,
		}
		neuronOrder[i] = id
		inverseNeuronOrder[id] = i
	}
	for i := numIn; i < numIn+numOut; i++ {
		id := counter.Next()
		neurons[id] = &Neuron{
			Type:       NeuronOutput,
			Activation: outActivation,
		}
		neuronOrder[i] = id
		inverseNeuronOrder[id] = i
	}
	return &Genotype{
		NumIn:              numIn,
		NumOut:             numOut,
		Neurons:            neurons,
		Synapses:           make(map[int]*Synapse),
		NeuronOrder:        neuronOrder,
		InverseNeuronOrder: inverseNeuronOrder,
	}
}

// Create a new Genotype which is a clone of `g`
func NewGenotypeCopy(g *Genotype) *Genotype {
	neurons := make(map[int]*Neuron)
	synapses := make(map[int]*Synapse)
	neuronOrder := make([]int, len(g.NeuronOrder))
	inverseNeuronOrder := make(map[int]int)
	copy(neuronOrder, g.NeuronOrder)
	for nid, n := range g.Neurons {
		neurons[nid] = &Neuron{n.Type, n.Activation}
	}
	for nid, no := range g.InverseNeuronOrder {
		inverseNeuronOrder[nid] = no
	}
	for sid, s := range g.Synapses {
		synapses[sid] = &Synapse{s.From, s.To, s.Weight}
	}
	return &Genotype{
		NumIn:              g.NumIn,
		NumOut:             g.NumOut,
		Neurons:            neurons,
		Synapses:           synapses,
		NeuronOrder:        neuronOrder,
		InverseNeuronOrder: inverseNeuronOrder,
	}
}

// Creates a new Genotype which is a clone of g1, but with 50% of the weights of matching synapses from g2.
// g1 should be the fitter parent
func NewGenotypeCrossover(g1 *Genotype, g2 *Genotype) *Genotype {
	neurons := make(map[int]*Neuron)
	synapses := make(map[int]*Synapse)
	neuronOrder := make([]int, len(g1.NeuronOrder))
	inverseNeuronOrder := make(map[int]int)
	// Copy neurons from g1. Currently all activations are copied too and not crossed over
	copy(neuronOrder, g1.NeuronOrder)
	for nid, n := range g1.Neurons {
		neurons[nid] = &Neuron{n.Type, n.Activation}
	}
	for nid, no := range g1.InverseNeuronOrder {
		inverseNeuronOrder[nid] = no
	}
	for sid, s1 := range g1.Synapses {
		s2, ok := g2.Synapses[sid]
		if ok {
			//50% chance to pick a weight from s2
			if rand.Intn(2) == 0 {
				synapses[sid] = &Synapse{s1.From, s1.To, s2.Weight}
			} else {
				synapses[sid] = &Synapse{s1.From, s1.To, s1.Weight}
			}
		} else {
			synapses[sid] = &Synapse{s1.From, s1.To, s1.Weight}
		}
	}
	return &Genotype{
		NumIn:              g1.NumIn,
		NumOut:             g1.NumOut,
		Neurons:            neurons,
		Synapses:           synapses,
		NeuronOrder:        neuronOrder,
		InverseNeuronOrder: inverseNeuronOrder,
	}
}

// Find the ID of the synapse that goes from `from` to `to`.
// Will return `synapseID, error`, where error might be caused by the connection not existing
func (n *Genotype) LookupSynapse(from, to int) (int, error) {
	for ci, c := range n.Synapses {
		if c.From == from && c.To == to {
			return ci, nil
		}
	}
	return -1, errors.New("cannot find connection")
}

// Create a synapse between two neurons using the provided Counter, with weight `weight`.
// `from` and `to` are the innovation IDs of two neurons.
// If `from` is ordered after `to` then the connection is recurrent.
// Will return `synapseID, error`
func (n *Genotype) AddSynapse(counter Counter, from, to int, weight float64) (int, error) {
	if !(n.IsNeuron(from) && n.IsNeuron(to)) {
		return -1, errors.New("ids are not both nodes")
	}
	if _, err := n.LookupSynapse(from, to); err == nil {
		return -1, errors.New("connection already exists")
	}
	nodeFrom := n.Neurons[from]
	nodeTo := n.Neurons[to]
	if nodeFrom.Type == NeuronOutput && nodeTo.Type == NeuronOutput {
		return -1, errors.New("cannot create connection from output to output")
	} else if nodeFrom.Type == NeuronInput && nodeTo.Type == NeuronInput {
		return -1, errors.New("cannot create connection from input to input")
	} else if from == to {
		return -1, errors.New("cannot connect neuron to itself")
	}
	/*nodeFromOrder, _ := n.GetNeuronOrder(from)
	nodeToOrder, _ := n.GetNeuronOrder(to)
	if nodeToOrder < nodeFromOrder {
		return -1, errors.New("recursive connections are not supported right now")
	}*/
	c := &Synapse{from, to, weight}
	id := counter.Next()
	n.Synapses[id] = c
	return id, nil
}

// Create a hidden neuron on a synapse using the provided Counter, with activation function `activation`.
// `conn` is the ID of the synapse to create the neuron on, and `conn` must not refer to a recurrent connection.
// Will return `neuronID, synapseID, error`, where synapseID is the id of the synapse that was created due to the original synapse being split.
// This synapse connects the new neuron to the old synapses endpoint, and the old synapses endpoint is moved to the new neuron
func (n *Genotype) AddNeuron(counter Counter, conn int, activation Activation) (int, int, error) {
	if !n.IsSynapse(conn) {
		return -1, -1, errors.New("id is not a connection")
	}
	// Get new ids
	newNodeID := counter.Next()
	newConnID := counter.Next()
	// Find info about existing conn
	connPtr := n.Synapses[conn]
	fromNode, toNode := connPtr.From, connPtr.To
	// Calculate and insert the order
	fromNodeOrder, _ := n.GetNeuronOrder(fromNode)
	toNodeOrder, _ := n.GetNeuronOrder(toNode)
	if toNodeOrder < fromNodeOrder {
		return -1, -1, errors.New("cannot create neuron on recursive connection")
	}
	var newNodeOrder int
	if fromNodeOrder < toNodeOrder {
		newNodeOrder = fromNodeOrder + 1
	} else {
		newNodeOrder = toNodeOrder + 1
	}
	if newNodeOrder < n.NumIn {
		newNodeOrder = n.NumIn
	}
	if newNodeOrder > len(n.NeuronOrder)-n.NumOut+1 {
		newNodeOrder = len(n.NeuronOrder) - n.NumOut + 1
	}
	connPtr.To = newNodeID
	n.insertNeuron(newNodeID, newNodeOrder, &Neuron{Type: NeuronHidden, Activation: activation})
	// Create a new conn from new node to the old endpoint
	newConn := &Synapse{newNodeID, toNode, 1}
	n.Synapses[newConnID] = newConn
	return newNodeID, newConnID, nil
}

// Prune the synapse with id `sid` from the genotype. Then recursively check to see if this has made any other synapses or neurons redundant, and remove those too.
func (n *Genotype) PruneSynapse(sid int) error {
	if !n.IsSynapse(sid) {
		return errors.New("not a synapse")
	}
	syn := n.Synapses[sid]
	from, to := syn.From, syn.To
	delete(n.Synapses, sid)
	if n.IsNeuron(from) && n.Neurons[from].Type == NeuronHidden && !n.hasNeuronOut(from) {
		n.removeNeuron(from)
		for sid, s := range n.Synapses {
			if s.To == from {
				n.PruneSynapse(sid)
			}
		}
	}
	if n.IsNeuron(to) && n.Neurons[to].Type == NeuronHidden && !n.hasNeuronIn(to) {
		n.removeNeuron(to)
		for sid, s := range n.Synapses {
			if s.From == to {
				n.PruneSynapse(sid)
			}
		}
	}
	return nil
}

func (n *Genotype) hasNeuronOut(nid int) bool {
	for _, s := range n.Synapses {
		if s.From == nid {
			return true
		}
	}
	return false
}
func (n *Genotype) hasNeuronIn(nid int) bool {
	for _, s := range n.Synapses {
		if s.To == nid {
			return true
		}
	}
	return false
}

// Return the number of input, hidden, and output neurons in this genotype
func (n *Genotype) Topology() (int, int, int) {
	return n.NumIn, len(n.NeuronOrder) - n.NumIn - n.NumOut, n.NumOut
}

// Check if a given id is a neuron in this genotype
func (n *Genotype) IsNeuron(id int) bool {
	_, ok := n.Neurons[id]
	return ok
}

// Check if a given id is a synapse in this genotype
func (n *Genotype) IsSynapse(id int) bool {
	_, ok := n.Synapses[id]
	return ok
}

// Gets the order (position in which the neurons are calculated) of a neuron ID
func (n *Genotype) GetNeuronOrder(nid int) (int, error) {
	if !n.IsNeuron(nid) {
		return -1, errors.New("not a node")
	}
	return n.InverseNeuronOrder[nid], nil
	/*for norder := range n.NeuronOrder {
		if n.NeuronOrder[norder] == nid {
			return norder, nil
		}
	}
	panic("node order list has become desynced. this should not have happened")*/
}

func insert[T any](a []T, index int, value T) []T {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func (g *Genotype) removeNeuron(nid int) {
	if g.Neurons[nid].Type != NeuronHidden {
		panic("oof")
	}
	o := g.InverseNeuronOrder[nid]
	delete(g.Neurons, nid)
	delete(g.InverseNeuronOrder, nid)
	g.NeuronOrder = append(g.NeuronOrder[:o], g.NeuronOrder[o+1:]...)
	for i := range g.NeuronOrder {
		g.InverseNeuronOrder[g.NeuronOrder[i]] = i
	}
}

func (g *Genotype) insertNeuron(nid, order int, n *Neuron) {
	g.Neurons[nid] = n
	g.NeuronOrder = insert(g.NeuronOrder, order, nid)
	for i := range g.NeuronOrder {
		g.InverseNeuronOrder[g.NeuronOrder[i]] = i
	}
}
