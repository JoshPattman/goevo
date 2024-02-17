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
	forwardSynapses       []SynapseID // With these three we just track which synapses are of what type
	backwardSynapses      []SynapseID // A synapse can NEVER change type
	selfSynapses          []SynapseID
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
	forwardSyanpses := make([]SynapseID, 0)
	backwardSyanpses := make([]SynapseID, 0)
	selfSyanpses := make([]SynapseID, 0)

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
		forwardSynapses:       forwardSyanpses,
		backwardSynapses:      backwardSyanpses,
		selfSynapses:          selfSyanpses,
	}
}

func (g *Genotype) AddRandomNeuron(counter *Counter, activations ...Activation) bool {
	if len(g.weights) == 0 {
		return false
	}

	widx := rand.Intn(len(g.forwardSynapses) + len(g.backwardSynapses)) // We should never add a weight on a self synapse
	var sid SynapseID
	if widx < len(g.forwardSynapses) {
		sid = g.forwardSynapses[widx]
	} else {
		sid = g.backwardSynapses[widx-len(g.forwardSynapses)]
	}

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

	// Find the two original endpoints orders, and also which was first and which was second
	ao, bo := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]

	// Add b to the index of its synapse class
	if ao < bo {
		g.forwardSynapses = append(g.forwardSynapses, newSid)
	} else if ao > bo {
		g.backwardSynapses = append(g.backwardSynapses, newSid)
	} else {
		g.selfSynapses = append(g.selfSynapses, newSid)
	}

	// Create a new neuron
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
		if ao == bo && !recurrent {
			continue // No self connections if non recurrent
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
		if !recurrent {
			g.forwardSynapses = append(g.forwardSynapses, sid)
		} else if ep.From == ep.To {
			g.selfSynapses = append(g.selfSynapses, sid)
		} else {
			g.backwardSynapses = append(g.backwardSynapses, sid)
		}
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
	ep := g.synapseEndpointLookup[sid]

	fo, to := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
	if fo < to {
		idx := slices.Index(g.forwardSynapses, sid)
		g.forwardSynapses[idx] = g.forwardSynapses[len(g.forwardSynapses)-1]
		g.forwardSynapses = g.forwardSynapses[:len(g.forwardSynapses)-1]
	} else if fo > to {
		idx := slices.Index(g.backwardSynapses, sid)
		g.backwardSynapses[idx] = g.backwardSynapses[len(g.backwardSynapses)-1]
		g.backwardSynapses = g.backwardSynapses[:len(g.backwardSynapses)-1]
	} else {
		idx := slices.Index(g.selfSynapses, sid)
		g.selfSynapses[idx] = g.selfSynapses[len(g.selfSynapses)-1]
		g.selfSynapses = g.selfSynapses[:len(g.selfSynapses)-1]
	}

	delete(g.weights, sid)
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

// Change the activation of a rnadom HIDDEN neuron to one of the supplied activations
func (g *Genotype) MutateRandomActivation(activations ...Activation) bool {
	numHidden := len(g.neuronOrder) - g.numInputs - g.numOutputs
	if numHidden <= 0 {
		return false
	}
	i := g.numInputs + rand.Intn(numHidden)
	g.activations[g.neuronOrder[i]] = activations[rand.Intn(len(activations))]
	return true
}

// Render this genotype to an image.Image, with a width and height in inches
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
		slices.Clone(g.forwardSynapses),
		slices.Clone(g.backwardSynapses),
		slices.Clone(g.selfSynapses),
	}

	return gc
}

// Simple crossover of the genotypes, where g is fitter than g2
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
		slices.Clone(g.forwardSynapses),
		slices.Clone(g.backwardSynapses),
		slices.Clone(g.selfSynapses),
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

func (g *Genotype) NumInputNeurons() int {
	return g.numInputs
}

func (g *Genotype) NumOutputNeurons() int {
	return g.numOutputs
}

func (g *Genotype) NumHiddenNeurons() int {
	return len(g.activations) - g.numInputs - g.numOutputs
}

func (g *Genotype) NumNeurons() int {
	return len(g.activations)
}

func (g *Genotype) NumSynapses() int {
	return len(g.weights)
}

// This will run as many checks as possible to check the genotype is valid.
// It is really only designed to be used as part of a test suite to catch errors with the package.
// This should never throw an error, but if it does either there is a bug in the package, or the user has somehow invalidated the genotype.
func (g *Genotype) Validate() error {
	// Check there are enough inputs and outputs
	if g.numInputs <= 0 {
		return fmt.Errorf("not enough inputs: %v", g.numInputs)
	}
	if g.numOutputs <= 0 {
		return fmt.Errorf("not enough outputs: %v", g.numOutputs)
	}

	// Check there are at least enough input and output nodes
	if len(g.neuronOrder) < g.numInputs+g.numOutputs {
		return fmt.Errorf("than number of inputs (%v) and outputs (%v) is not possible with the number of nodes loaded (%v)", g.numInputs, g.numOutputs, len(g.neuronOrder))
	}

	// Check max synapse value is valid
	if g.maxSynapseValue <= 0 {
		return fmt.Errorf("invalid maximum synapse value: %v", g.maxSynapseValue)
	}

	// Ensure that all node indexes have same length
	if len(g.neuronOrder) != len(g.inverseNeuronOrder) {
		return fmt.Errorf("inverse neuron order has length %v but neuron order has length %v", len(g.inverseNeuronOrder), len(g.neuronOrder))
	}
	if len(g.neuronOrder) != len(g.activations) {
		return fmt.Errorf("activations has length %v but neuron order has length %v", len(g.activations), len(g.neuronOrder))
	}

	// Ensure that all weight indexes have same length
	if len(g.weights) != len(g.synapseEndpointLookup) {
		return fmt.Errorf("synapse endpoint lookup has length %v but weights has length %v", len(g.synapseEndpointLookup), len(g.weights))
	}
	if len(g.weights) != len(g.endpointSynapseLookup) {
		return fmt.Errorf("endpoint synapse lookup has length %v but weights has length %v", len(g.endpointSynapseLookup), len(g.weights))
	}
	if len(g.weights) != len(g.forwardSynapses)+len(g.backwardSynapses)+len(g.selfSynapses) {
		return fmt.Errorf("forward, backward, and self synapses have combined length %v but weights has length %v", len(g.forwardSynapses)+len(g.backwardSynapses)+len(g.selfSynapses), len(g.weights))
	}

	// Ensure that there are no ids that are the same between the neurons and the synapses
	foundIDs := make(map[int]bool)
	for id := range g.activations {
		if _, ok := foundIDs[int(id)]; ok {
			return fmt.Errorf("repeated id: %v", id)
		}
		foundIDs[int(id)] = true
	}
	for id := range g.weights {
		if _, ok := foundIDs[int(id)]; ok {
			return fmt.Errorf("repeated id: %v", id)
		}
		foundIDs[int(id)] = true
	}

	// Check that synapseEPLookup and EPSynapseLookup are synced
	for id, ep := range g.synapseEndpointLookup {
		if id2, ok := g.endpointSynapseLookup[ep]; !ok {
			return fmt.Errorf("missing id that exists in synapse endpoint lookup but not in endpoint synapse lookup: %v", id)
		} else if id != id2 {
			return fmt.Errorf("synapse endpoint lookup and endpoint synapse lookup are not symmetrical with id: %v (there and back becomes %v)", id, id2)
		}
	}
	for ep, id := range g.endpointSynapseLookup {
		if ep2, ok := g.synapseEndpointLookup[id]; !ok {
			return fmt.Errorf("missing id that exists in endpoint synapse lookup but not in synapse endpoint lookup: %v", id)
		} else if ep != ep2 {
			return fmt.Errorf("synapse endpoint lookup and endpoint synapse lookup are not symmetrical with ep: %v (there and back becomes %v)", ep, ep2)
		}
	}

	return nil
}
