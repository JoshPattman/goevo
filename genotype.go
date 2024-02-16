package goevo

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"slices"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"golang.org/x/exp/maps"
)

type NeuronID int
type SynapseID int

// SynapseEP is the endpoints of a synapse
type SynapseEP struct {
	From NeuronID
	To   NeuronID
}

// Genotype represents the DNA of a creature. It is optimised for mutating, but cannot be run directly.
type Genotype struct {
	maxSynapseValue       float64
	numInputs             int
	numOutputs            int
	neuronOrder           []NeuronID
	inverseNeuronOrder    map[NeuronID]int
	activations           map[NeuronID]Activation
	weights               map[SynapseID]float64
	synapseEndpointLookup map[SynapseID]SynapseEP
	endpointSynapseLookup map[SynapseEP]SynapseID
}

func NewGenotype(counter *Counter, inputs, outputs int, outputActivation Activation) *Genotype {
	if inputs <= 0 || outputs <= 0 {
		panic("must have at least one input and one output")
	}
	neuronOrder := make([]NeuronID, 0)
	inverseNeuronOrder := make(map[NeuronID]int)
	activations := make(map[NeuronID]Activation)
	weights := make(map[SynapseID]float64)
	synapseEndpointLookup := make(map[SynapseID]SynapseEP)
	endpointSynapseLookup := make(map[SynapseEP]SynapseID)

	for i := 0; i < inputs; i++ {
		id := NeuronID(counter.Next())
		neuronOrder = append(neuronOrder, id)
		inverseNeuronOrder[id] = len(neuronOrder) - 1
		activations[id] = Linear
	}

	for i := 0; i < outputs; i++ {
		id := NeuronID(counter.Next())
		neuronOrder = append(neuronOrder, id)
		inverseNeuronOrder[id] = len(neuronOrder) - 1
		activations[id] = outputActivation
	}

	return &Genotype{
		maxSynapseValue:       3,
		numInputs:             inputs,
		numOutputs:            outputs,
		neuronOrder:           neuronOrder,
		inverseNeuronOrder:    inverseNeuronOrder,
		activations:           activations,
		weights:               weights,
		synapseEndpointLookup: synapseEndpointLookup,
		endpointSynapseLookup: endpointSynapseLookup,
	}
}

func (g *Genotype) AddRandomNeuron(counter *Counter, activations ...Activation) bool {
	if len(g.weights) == 0 {
		return false
	}

	sid := randomMapKey(g.weights)

	ep := g.synapseEndpointLookup[sid]

	newSid := SynapseID(counter.Next())
	newNid := NeuronID(counter.Next())

	epa := SynapseEP{ep.From, newNid}
	epb := SynapseEP{newNid, ep.To}

	// Swap the old connection for a, which will also retain the original weight
	delete(g.endpointSynapseLookup, ep)
	g.endpointSynapseLookup[epa] = sid
	g.synapseEndpointLookup[sid] = epa
	// Don't need to modify weights because weights[sid] is already there

	// Create a new connection for b, with weight of 1 to minimise affect on behaviour
	g.endpointSynapseLookup[epb] = newSid
	g.synapseEndpointLookup[newSid] = epb
	g.weights[newSid] = 1

	// Create a new neuron
	// Find the two original endpoints orders, and also which was first and which was second
	ao, bo := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
	firstO, secondO := ao, bo
	if bo < ao {
		firstO, secondO = bo, ao
	}

	// Check that the synapse is valid. If it is not, somthing has gone wrong
	if g.isInputOrder(ao) && g.isInputOrder(bo) {
		panic("trying to insert a node on a connection between two inputs, either a bug has occured or you have created an invalid genotype")
	} else if g.isOutputOrder(ao) && g.isOutputOrder(bo) {
		panic("trying to insert a node on a connection between two outputs, either a bug has occured or you have created an invalid genotype")
	} else if g.isInputOrder(bo) {
		panic("trying to insert a node on a connection that ends in an input, either a bug has occured or you have created an invalid genotype")
	}
	// Find the new order of the neuron
	no := int(math.Round((float64(ao) + float64(bo)) / 2.0))     // Find the order that is halfway between them
	startPosition := max(g.numInputs, firstO+1)                  // First valid position INCLUSIVE
	endPosition := min(len(g.neuronOrder)-g.numOutputs, secondO) // Last valid position INCLUSIVE
	if startPosition > endPosition {
		panic("failed to find valid placement of neuron, this should not have happened")
	}
	if no < startPosition {
		no = startPosition
	} else if no > endPosition {
		no = endPosition
	}

	// Insert the neuron at that order
	newNeuronOrder := make([]NeuronID, len(g.neuronOrder)+1)
	copy(newNeuronOrder, g.neuronOrder[:no])
	newNeuronOrder[no] = newNid
	copy(newNeuronOrder[no+1:], g.neuronOrder[no:])
	g.neuronOrder = newNeuronOrder
	for i := no; i < len(g.neuronOrder); i++ {
		g.inverseNeuronOrder[g.neuronOrder[i]] = i
	}
	// Add the activation
	g.activations[newNid] = activations[rand.Intn(len(activations))]

	return true
}

func (g *Genotype) AddRandomSynapse(counter *Counter, weightStd float64, recurrent bool) bool {
	// Almost always find a new connection after 10 tries
	for i := 0; i < 10; i++ {
		ao := rand.Intn(len(g.neuronOrder))
		bo := rand.Intn(len(g.neuronOrder))
		if ao == bo {
			continue // No self connections
		}
		if (!recurrent && ao > bo) || (recurrent && bo > ao) {
			ao, bo = bo, ao // Ensure that this connection is of correct type
		}
		if (g.isInputOrder(bo)) || (g.isOutputOrder(ao) && g.isOutputOrder(bo)) {
			continue // Trying to connect either anything-input or output-output
		}
		aid, bid := g.neuronOrder[ao], g.neuronOrder[bo]
		ep := SynapseEP{aid, bid}
		if _, ok := g.endpointSynapseLookup[ep]; ok {
			continue // This connection already exists, try to find another
		}
		sid := SynapseID(counter.Next())
		g.endpointSynapseLookup[ep] = sid
		g.synapseEndpointLookup[sid] = ep
		g.weights[sid] = clamp(rand.NormFloat64()*weightStd, -g.maxSynapseValue, g.maxSynapseValue)
		return true
	}
	return false
}

func (g *Genotype) MutateRandomSynapse(std float64) bool {
	if len(g.weights) == 0 {
		return false
	}

	sid := randomMapKey(g.weights)
	g.weights[sid] = clamp(g.weights[sid]+rand.NormFloat64()*std, -g.maxSynapseValue, g.maxSynapseValue)

	return true
}

// This will delete a random synapse. It will leave hanging neurons, because they may be useful later.
func (g *Genotype) RemoveRandomSynapse() bool {
	if len(g.weights) == 0 {
		return false
	}
	sid := randomMapKey(g.weights)
	delete(g.weights, sid)
	ep := g.synapseEndpointLookup[sid]
	delete(g.synapseEndpointLookup, sid)
	delete(g.endpointSynapseLookup, ep)
	return true
}

// This will set the weight of a random synapse to 0. Kind of similar to disabling a synapse, which this implementation does not have.
func (g *Genotype) ResetRandomSynapse() bool {
	if len(g.weights) == 0 {
		return false
	}
	sid := randomMapKey(g.weights)
	g.weights[sid] = 0
	return true
}

func (g *Genotype) MutateRandomActivation(activations ...Activation) bool {
	numHidden := len(g.neuronOrder) - g.numInputs - g.numOutputs
	if numHidden <= 0 {
		return false
	}
	i := g.numInputs + rand.Intn(numHidden)
	g.activations[g.neuronOrder[i]] = activations[rand.Intn(len(activations))]
	return true
}

func (g *Genotype) Draw(width, height float64) image.Image {
	gv := graphviz.New()

	graph, err := gv.Graph()
	if err != nil {
		panic(fmt.Sprintf("error when creating a graph, this should not have happened: %v", err))
	}

	defer func() {
		graph.Close()
		gv.Close()
	}()

	graph.SetRankDir(cgraph.LRRank)
	graph.SetRatio(cgraph.FillRatio)
	graph.SetSize(width, height)

	nodes := make(map[NeuronID]*cgraph.Node)
	for no, nid := range g.neuronOrder {
		nodes[nid], err = graph.CreateNode(fmt.Sprintf("N%v [%v]\n%v", nid, no, g.activations[nid]))
		if err != nil {
			panic(fmt.Sprintf("error when creating node on a graph, this should not have happened: %v", err))
		}
		if no < g.numInputs {
			nodes[nid].SetColor("green")
		} else if no >= len(g.neuronOrder)-g.numOutputs {
			nodes[nid].SetColor("red")
		}
		nodes[nid].SetShape(cgraph.RectShape)
	}

	for wid, w := range g.weights {
		ep := g.synapseEndpointLookup[wid]
		edge, _ := graph.CreateEdge(fmt.Sprintf("%v->%v", ep.From, ep.To), nodes[ep.From], nodes[ep.To])
		edge.SetLabel(fmt.Sprintf("%.3f", w))
	}
	img, err := gv.RenderImage(graph)
	if err != nil {
		panic(fmt.Sprintf("error when creating an image, this should not have happened: %v", err))
	}
	return img
}

// clones the genotype
func (g *Genotype) Clone() *Genotype {
	gc := &Genotype{
		g.maxSynapseValue,
		g.numInputs,
		g.numOutputs,
		slices.Clone(g.neuronOrder),
		maps.Clone(g.inverseNeuronOrder),
		maps.Clone(g.activations),
		maps.Clone(g.weights),
		maps.Clone(g.synapseEndpointLookup),
		maps.Clone(g.endpointSynapseLookup),
	}

	return gc
}

// g is fitter than g2
func (g *Genotype) CrossoverWith(g2 *Genotype) *Genotype {
	gc := &Genotype{
		g.maxSynapseValue,
		g.numInputs,
		g.numOutputs,
		slices.Clone(g.neuronOrder),
		maps.Clone(g.inverseNeuronOrder),
		maps.Clone(g.activations),
		maps.Clone(g.weights),
		maps.Clone(g.synapseEndpointLookup),
		maps.Clone(g.endpointSynapseLookup),
	}

	for sid, sw := range g2.weights {
		if _, ok := gc.weights[sid]; ok {
			if rand.Float64() > 0.5 {
				gc.weights[sid] = sw
			}
		}
	}

	return gc
}

func (g *Genotype) isInputOrder(order int) bool {
	return order < g.numInputs
}

func (g *Genotype) isOutputOrder(order int) bool {
	return order >= len(g.neuronOrder)-g.numOutputs
}

func (g *Genotype) NumInputs() int {
	return g.numInputs
}

func (g *Genotype) NumOutputs() int {
	return g.numOutputs
}

func (g *Genotype) NumHidden() int {
	return len(g.activations) - g.numInputs - g.numOutputs
}

func (g *Genotype) NumNeurons() int {
	return len(g.activations)
}

func (g *Genotype) NumSynapses() int {
	return len(g.weights)
}
