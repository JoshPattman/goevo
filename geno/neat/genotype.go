// Package neat provides a genotype for a neural network using the NEAT algorithm.
// This does not provide the NEAT population, that is in pop/neatpop
package neat

import (
	"fmt"
	"math"
	"math/rand"
	"slices"

	"github.com/JoshPattman/goevo"
	"golang.org/x/exp/maps"
)

var _ goevo.Cloneable = &Genotype{}
var _ goevo.PointCrossoverable = &Genotype{}

// NeuronID is the unique identifier for a neuron in a NEATGenotype
type NeuronID int

// SynapseID is the unique identifier for a synapse in a NEATGenotype
type SynapseID int

// SynapseEP is the endpoints of a synapse in a NEATGenotype
type SynapseEP struct {
	From NeuronID
	To   NeuronID
}

// Genotype is a genotype for a neural network using the NEAT algorithm.
// It is conceptually similar to the DNA of an organism: it encodes how to build a neural network, but is not the neural network itself.
// This means if you want to actually run the neural network, you need to use the [Genotype.Build] method to create a [Phenotype].
type Genotype struct {
	maxSynapseValue       float64
	numInputs             int
	numOutputs            int
	neuronOrder           []NeuronID
	inverseNeuronOrder    map[NeuronID]int
	activations           map[NeuronID]goevo.Activation
	weights               map[SynapseID]float64
	synapseEndpointLookup map[SynapseID]SynapseEP
	endpointSynapseLookup map[SynapseEP]SynapseID
	forwardSynapses       []SynapseID // With these three we just track which synapses are of what type
	backwardSynapses      []SynapseID // A synapse can NEVER change type
	selfSynapses          []SynapseID
}

// NewGenotype creates a new NEATGenotype with the given number of inputs and outputs, and the given output activation function.
// All output neurons will have the same activation function, and all input neurons will have the linear activation function.
// The genotype will have no synapses.
func NewGenotype(counter *goevo.Counter, inputs, outputs int, outputActivation goevo.Activation) *Genotype {
	if inputs <= 0 || outputs <= 0 {
		panic("must have at least one input and one output")
	}
	neuronOrder := make([]NeuronID, 0)
	inverseNeuronOrder := make(map[NeuronID]int)
	activations := make(map[NeuronID]goevo.Activation)
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
		activations[id] = goevo.Linear
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

// AddRandomNeuron adds a new neuron to the genotype on a random forward synapse.
// It will return false if there are no forward synapses to add to.
// The new neuron will have a random activation function from the given list of activations.
func (g *Genotype) AddRandomNeuron(counter *goevo.Counter, activations ...goevo.Activation) bool {
	if len(g.forwardSynapses) == 0 {
		return false
	}

	// We only ever want to add nodes on forward synapses
	sid := g.forwardSynapses[rand.Intn(len(g.forwardSynapses))]

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

// AddRandomSynapse adds a new synapse to the genotype between two nodes.
// It will return false if it failed to find a place to put the synapse after 10 tries.
// The synapse will have a random weight from a normal distribution with the given standard deviation.
// If recurrent is true, the synapse will be recurrent, otherwise it will not.
func (g *Genotype) AddRandomSynapse(counter *goevo.Counter, weightStd float64, recurrent bool) bool {
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

// MutateRandomSynapse will change the weight of a random synapse by a random amount from a normal distribution with the given standard deviation.
// It will return false if there are no synapses to mutate.
func (g *Genotype) MutateRandomSynapse(std float64) bool {
	if len(g.weights) == 0 {
		return false
	}

	sid := randomMapKey(g.weights)
	g.weights[sid] = clamp(g.weights[sid]+rand.NormFloat64()*std, -g.maxSynapseValue, g.maxSynapseValue)

	return true
}

// RemoveRandomSynapse will remove a random synapse from the genotype.
// It will return false if there are no synapses to remove.
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

// ResetRandomSynapse will reset the weight of a random synapse to 0.
// It will return false if there are no synapses to reset.
func (g *Genotype) ResetRandomSynapse() bool {
	if len(g.weights) == 0 {
		return false
	}
	sid := randomMapKey(g.weights)
	g.weights[sid] = 0
	return true
}

// MutateRandomActivation will change the activation function of a random hidden neuron to
// a random activation function from the given list of activations.
// It will return false if there are no hidden neurons to mutate.
func (g *Genotype) MutateRandomActivation(activations ...goevo.Activation) bool {
	numHidden := len(g.neuronOrder) - g.numInputs - g.numOutputs
	if numHidden <= 0 {
		return false
	}
	i := g.numInputs + rand.Intn(numHidden)
	g.activations[g.neuronOrder[i]] = activations[rand.Intn(len(activations))]
	return true
}

// Clone returns a new genotype that is an exact copy of this genotype.
func (g *Genotype) Clone() any {
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

// PointCrossoverWith will return a new genotype that is a crossover of this genotype with the given genotype.
// This crossover is similar to what is used in the original NEAT implementation,
// where only the weights of the synapses are crossed over (the entire structure of g is kept the same).
// For this reason, the first genotype g should be the fitter parent.
func (g *Genotype) PointCrossoverWith(g2i goevo.PointCrossoverable) goevo.PointCrossoverable {
	g2 := g2i.(*Genotype)
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

// NumInputNeurons returns the number of input neurons in the genotype.
func (g *Genotype) NumInputNeurons() int {
	return g.numInputs
}

// NumOutputNeurons returns the number of output neurons in the genotype.
func (g *Genotype) NumOutputNeurons() int {
	return g.numOutputs
}

// NumHiddenNeurons returns the number of hidden neurons in the genotype.
func (g *Genotype) NumHiddenNeurons() int {
	return len(g.activations) - g.numInputs - g.numOutputs
}

// NumNeurons returns the total number of neurons in the genotype.
func (g *Genotype) NumNeurons() int {
	return len(g.activations)
}

// NumSynapses returns the total number of synapses in the genotype.
func (g *Genotype) NumSynapses() int {
	return len(g.weights)
}

// Validate runs as many checks as possible to check the genotype is valid.
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

	// Check that synapseEPLookup and EPSynapseLookup are synced.
	// Only need to do this one way because we have already checked that they have same length
	for id, ep := range g.synapseEndpointLookup {
		if id2, ok := g.endpointSynapseLookup[ep]; !ok {
			return fmt.Errorf("missing id that exists in synapse endpoint lookup but not in endpoint synapse lookup: %v", id)
		} else if id != id2 {
			return fmt.Errorf("synapse endpoint lookup and endpoint synapse lookup are not symmetrical with id: %v (there and back becomes %v)", id, id2)
		}
	}

	// Check that weights and synapseEPLookup are synced.
	// Again, they already have same length.
	for id := range g.synapseEndpointLookup {
		if _, ok := g.weights[id]; !ok {
			return fmt.Errorf("missing id that exists in synapse endpoint lookup but not in weights: %v", id)
		}
	}

	// Check that neuron order and inverse neuron order are synced
	for i := range g.neuronOrder {
		if g.inverseNeuronOrder[g.neuronOrder[i]] != i {
			return fmt.Errorf("order %v is not symmetrical in neuron order and inverse neuron order", i)
		}
	}

	// Check that all classes of synapse are correctly categorised
	for _, sid := range g.forwardSynapses {
		ep := g.synapseEndpointLookup[sid]
		of, ot := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
		if ot <= of {
			return fmt.Errorf("synapse with id %v is incorrectly categorised as forward", sid)
		}
	}
	for _, sid := range g.backwardSynapses {
		ep := g.synapseEndpointLookup[sid]
		of, ot := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
		if ot >= of {
			return fmt.Errorf("synapse with id %v is incorrectly categorised as backward", sid)
		}
	}
	for _, sid := range g.selfSynapses {
		ep := g.synapseEndpointLookup[sid]
		of, ot := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
		if ot != of {
			return fmt.Errorf("synapse with id %v is incorrectly categorised as self", sid)
		}
	}

	return nil
}
