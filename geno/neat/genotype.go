// Package neat provides a genotype for a neural network using the NEAT algorithm.
// This does not provide the NEAT population, that is in pop/neatpop
package neat

import (
	"fmt"
	"slices"

	"github.com/JoshPattman/goevo"
	"golang.org/x/exp/maps"
)

var _ goevo.Cloneable = &Genotype{}
var _ goevo.Buildable = &Genotype{}

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
