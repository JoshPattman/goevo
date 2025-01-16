package goevo

import (
	"encoding/json"
	"fmt"
	"image"
	"maps"
	"math"
	"math/rand"
	"slices"

	"github.com/goccy/go-graphviz"
)

// NeatNeuronID is the unique identifier for a neuron in a NEATGenotype
type NeatNeuronID int

// NeatSynapseID is the unique identifier for a synapse in a NEATGenotype
type NeatSynapseID int

// NeatSynapseEP is the endpoints of a synapse in a NEATGenotype
type NeatSynapseEP struct {
	From NeatNeuronID
	To   NeatNeuronID
}

// NeatGenotype is a genotype for a neural network using the NEAT algorithm.
// It is conceptually similar to the DNA of an organism: it encodes how to build a neural network, but is not the neural network itself.
// This means if you want to actually run the neural network, you need to use the [NeatGenotype.Build] method to create a [NeatPhenotype].
type NeatGenotype struct {
	maxSynapseValue       float64
	numInputs             int
	numOutputs            int
	neuronOrder           []NeatNeuronID
	inverseNeuronOrder    map[NeatNeuronID]int
	activations           map[NeatNeuronID]Activation
	weights               map[NeatSynapseID]float64
	synapseEndpointLookup map[NeatSynapseID]NeatSynapseEP
	endpointSynapseLookup map[NeatSynapseEP]NeatSynapseID
	forwardSynapses       []NeatSynapseID // With these three we just track which synapses are of what type
	backwardSynapses      []NeatSynapseID // A synapse can NEVER change type
	selfSynapses          []NeatSynapseID
}

// NewNeatGenotype creates a new NEATGenotype with the given number of inputs and outputs, and the given output activation function.
// All output neurons will have the same activation function, and all input neurons will have the linear activation function.
// The genotype will have no synapses.
func NewNeatGenotype(counter *Counter, inputs, outputs int, outputActivation Activation) *NeatGenotype {
	if inputs <= 0 || outputs <= 0 {
		panic("must have at least one input and one output")
	}
	neuronOrder := make([]NeatNeuronID, 0)
	inverseNeuronOrder := make(map[NeatNeuronID]int)
	activations := make(map[NeatNeuronID]Activation)
	weights := make(map[NeatSynapseID]float64)
	synapseEndpointLookup := make(map[NeatSynapseID]NeatSynapseEP)
	endpointSynapseLookup := make(map[NeatSynapseEP]NeatSynapseID)
	forwardSyanpses := make([]NeatSynapseID, 0)
	backwardSyanpses := make([]NeatSynapseID, 0)
	selfSyanpses := make([]NeatSynapseID, 0)

	for i := 0; i < inputs; i++ {
		id := NeatNeuronID(counter.Next())
		neuronOrder = append(neuronOrder, id)
		inverseNeuronOrder[id] = len(neuronOrder) - 1
		activations[id] = Linear
	}

	for i := 0; i < outputs; i++ {
		id := NeatNeuronID(counter.Next())
		neuronOrder = append(neuronOrder, id)
		inverseNeuronOrder[id] = len(neuronOrder) - 1
		activations[id] = outputActivation
	}

	return &NeatGenotype{
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
func (g *NeatGenotype) Clone() any {
	gc := &NeatGenotype{
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

func (g *NeatGenotype) isInputOrder(order int) bool {
	return order < g.numInputs
}

func (g *NeatGenotype) isOutputOrder(order int) bool {
	return order >= len(g.neuronOrder)-g.numOutputs
}

// NumInputNeurons returns the number of input neurons in the genotype.
func (g *NeatGenotype) NumInputNeurons() int {
	return g.numInputs
}

// NumOutputNeurons returns the number of output neurons in the genotype.
func (g *NeatGenotype) NumOutputNeurons() int {
	return g.numOutputs
}

// NumHiddenNeurons returns the number of hidden neurons in the genotype.
func (g *NeatGenotype) NumHiddenNeurons() int {
	return len(g.activations) - g.numInputs - g.numOutputs
}

// NumNeurons returns the total number of neurons in the genotype.
func (g *NeatGenotype) NumNeurons() int {
	return len(g.activations)
}

// NumSynapses returns the total number of synapses in the genotype.
func (g *NeatGenotype) NumSynapses() int {
	return len(g.weights)
}

// Validate runs as many checks as possible to check the genotype is valid.
// It is really only designed to be used as part of a test suite to catch errors with the package.
// This should never throw an error, but if it does either there is a bug in the package, or the user has somehow invalidated the genotype.
func (g *NeatGenotype) Validate() error {
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

// AddRandomNeuron adds a new neuron to the genotype on a random forward synapse.
// It will return false if there are no forward synapses to add to.
// The new neuron will have a random activation function from the given list of activations.
func (g *NeatGenotype) AddRandomNeuron(counter *Counter, activations ...Activation) bool {
	if len(g.forwardSynapses) == 0 {
		return false
	}

	// We only ever want to add nodes on forward synapses
	sid := g.forwardSynapses[rand.Intn(len(g.forwardSynapses))]

	ep := g.synapseEndpointLookup[sid]

	newSid := NeatSynapseID(counter.Next())
	newNid := NeatNeuronID(counter.Next())

	epa := NeatSynapseEP{ep.From, newNid}
	epb := NeatSynapseEP{newNid, ep.To}

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
	newNeuronOrder := make([]NeatNeuronID, len(g.neuronOrder)+1)
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
func (g *NeatGenotype) AddRandomSynapse(counter *Counter, weightStd float64, recurrent bool) bool {
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
		ep := NeatSynapseEP{aid, bid}
		if _, ok := g.endpointSynapseLookup[ep]; ok {
			continue // This connection already exists, try to find another
		}
		sid := NeatSynapseID(counter.Next())
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
func (g *NeatGenotype) MutateRandomSynapse(std float64) bool {
	if len(g.weights) == 0 {
		return false
	}

	sid := randomMapKey(g.weights)
	g.weights[sid] = clamp(g.weights[sid]+rand.NormFloat64()*std, -g.maxSynapseValue, g.maxSynapseValue)

	return true
}

// RemoveRandomSynapse will remove a random synapse from the genotype.
// It will return false if there are no synapses to remove.
func (g *NeatGenotype) RemoveRandomSynapse() bool {
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
func (g *NeatGenotype) ResetRandomSynapse() bool {
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
func (g *NeatGenotype) MutateRandomActivation(activations ...Activation) bool {
	numHidden := len(g.neuronOrder) - g.numInputs - g.numOutputs
	if numHidden <= 0 {
		return false
	}
	i := g.numInputs + rand.Intn(numHidden)
	g.activations[g.neuronOrder[i]] = activations[rand.Intn(len(activations))]
	return true
}

// neatMutationStd is a reproduction strategy that uses a standard deviation for the number of mutations in each category.
// The standard deviation is not scaled by the size of the network, meaning that larger networks will tend to have more mutations than smaller networks.
type neatMutationStd struct {
	// The standard deviation for the number of new synapses
	stdNumNewSynapses float64
	// The standard deviation for the number of new recurrent synapses
	stdNumNewRecurrentSynapses float64
	// The standard deviation for the number of new neurons
	stdNumNewNeurons float64
	// The standard deviation for the number of synapses to mutate
	stdNumMutateSynapses float64
	// The standard deviation for the number of synapses to prune
	stdNumPruneSynapses float64
	// The standard deviation for the number of activations to mutate
	stdNumMutateActivations float64

	// The standard deviation for the weight of new synapses
	stdNewSynapseWeight float64
	// The standard deviation for the weight of mutated synapses
	stdMutateSynapseWeight float64

	// The maximum number of hidden neurons this mutation can add
	maxHiddenNeurons int

	// The counter to use for new synapse IDs
	counter *Counter
	// The possible activations to use for new neurons
	possibleActivations []Activation
}

func NewNeatMutationStd(
	counter *Counter,
	activations []Activation,
	stdNumNewForwardSynapses float64,
	stdNumNewRecurrentSynapses float64,
	stdNumNewNeurons float64,
	stdNumMutateSynapses float64,
	stdNumPruneSynapses float64,
	stdNumMutateActivations float64,
	stdNewSynapseWeight float64,
	stdMutateSynapseWeight float64,
	maxHiddenNeurons int,
) Mutation[*NeatGenotype] {
	if counter == nil {
		panic("cannot have nil counter")
	}
	if len(activations) == 0 {
		panic("cannot have no activations")
	}
	// TODO: Should probably check stds are all above 0 but it wont break anything
	return &neatMutationStd{
		counter:                    counter,
		possibleActivations:        activations,
		stdNumNewSynapses:          stdNumNewForwardSynapses,
		stdNumNewRecurrentSynapses: stdNumNewRecurrentSynapses,
		stdNumNewNeurons:           stdNumNewNeurons,
		stdNumMutateSynapses:       stdNumMutateSynapses,
		stdNumPruneSynapses:        stdNumPruneSynapses,
		stdNumMutateActivations:    stdNumMutateActivations,
		stdNewSynapseWeight:        stdNewSynapseWeight,
		stdMutateSynapseWeight:     stdMutateSynapseWeight,
		maxHiddenNeurons:           maxHiddenNeurons,
	}
}

// Reproduce creates a new genotype by crossing over and mutating the given genotypes.
func (r *neatMutationStd) Mutate(g *NeatGenotype) {
	for i := 0; i < stdN(r.stdNewSynapseWeight); i++ {
		g.AddRandomSynapse(r.counter, r.stdNewSynapseWeight, false)
	}
	for i := 0; i < stdN(r.stdNumNewRecurrentSynapses); i++ {
		g.AddRandomSynapse(r.counter, r.stdNewSynapseWeight, true)
	}
	for i := 0; i < stdN(r.stdNumNewNeurons); i++ {
		if r.maxHiddenNeurons < 0 || g.NumHiddenNeurons() < r.maxHiddenNeurons {
			g.AddRandomNeuron(r.counter, r.possibleActivations...)
		}
	}
	for i := 0; i < stdN(r.stdNumMutateSynapses); i++ {
		g.MutateRandomSynapse(r.stdMutateSynapseWeight)
	}
	for i := 0; i < stdN(r.stdNumPruneSynapses); i++ {
		g.RemoveRandomSynapse()
	}
	for i := 0; i < stdN(r.stdNumMutateActivations); i++ {
		g.MutateRandomActivation(r.possibleActivations...)
	}
}

type neatCrossoverSimple struct{}

func NewNeatCrossoverSimple() Crossover[*NeatGenotype] {
	return &neatCrossoverSimple{}
}

// Crossover implements CrossoverStrategy.
func (s *neatCrossoverSimple) Crossover(gs []*NeatGenotype) *NeatGenotype {
	if len(gs) != 2 {
		panic("expected 2 parents for simple crossover")
	}
	g, g2 := gs[0], gs[1]
	gc := &NeatGenotype{
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

// NumParents implements CrossoverStrategy.
func (s *neatCrossoverSimple) NumParents() int {
	return 2
}

type neatCrossoverAsexual struct{}

func NewNeatCrossoverAsexual() Crossover[*NeatGenotype] {
	return &neatCrossoverAsexual{}
}

func (s *neatCrossoverAsexual) Crossover(gs []*NeatGenotype) *NeatGenotype {
	if len(gs) != 1 {
		panic("expected 1 parent for asexual crossover")
	}
	return Clone(gs[0])
}

func (s *neatCrossoverAsexual) NumParents() int {
	return 1
}

type phenotypeConnection struct {
	toIdx int
	w     float64
}

// NeatPhenotype is a phenotype for a NEAT genotype.
// It conceptually represents a neural network, built according to the instructions in the NEATGenotype (DNA).
// Once built, the NeatPhenotype can be used to forward propagate inputs through the network,
// but it cannot be modified though mutation or corss-over.
type NeatPhenotype struct {
	numIn            int
	numOut           int
	accumulators     []float64
	lastAccumulators []float64
	activations      []Activation
	forwardWeights   [][]phenotypeConnection
	recurrentWeights [][]phenotypeConnection
}

// Build a NEATPhenotype from a NEATGenotype.
func (g *NeatGenotype) Build() Forwarder {
	accs := make([]float64, len(g.neuronOrder))
	laccs := make([]float64, len(g.neuronOrder))
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
	return &NeatPhenotype{
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
func (p *NeatPhenotype) Forward(x []float64) []float64 {
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
		p.accumulators[i] = Activate(p.accumulators[i], p.activations[i])
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
func (p *NeatPhenotype) Reset() {
	for i := range p.accumulators {
		p.accumulators[i] = 0
	}
}

// Make sure we implement json marshalling
var _ json.Marshaler = &NeatGenotype{}
var _ json.Unmarshaler = &NeatGenotype{}

type marshallableNeuron struct {
	ID         NeatNeuronID `json:"id"`
	Activation Activation   `json:"activation"`
}

type marshallableSynapse struct {
	ID     NeatSynapseID `json:"id"`
	From   NeatNeuronID  `json:"from"`
	To     NeatNeuronID  `json:"to"`
	Weight float64       `json:"weight"`
}

type marshallableGenotype struct {
	NumIn         int                   `json:"num_in"`
	NumOut        int                   `json:"num_out"`
	Neurons       []marshallableNeuron  `json:"neurons"`
	Synapses      []marshallableSynapse `json:"synapses"`
	MaxSynapseVal float64               `json:"max_synapse_val"`
}

// MarshalJSON implements json.Marshaler, allowing the genotype to be marshalled to JSON.
func (g *NeatGenotype) MarshalJSON() ([]byte, error) {
	mns := make([]marshallableNeuron, len(g.neuronOrder))
	for no, nid := range g.neuronOrder {
		mns[no] = marshallableNeuron{nid, g.activations[nid]}
	}
	mss := make([]marshallableSynapse, 0, len(g.weights))
	for sid, w := range g.weights {
		mss = append(mss, marshallableSynapse{
			ID:     sid,
			From:   g.synapseEndpointLookup[sid].From,
			To:     g.synapseEndpointLookup[sid].To,
			Weight: w,
		})
	}
	mg := marshallableGenotype{g.numInputs, g.numOutputs, mns, mss, g.maxSynapseValue}
	return json.Marshal(&mg)
}

// UnmarshalJSON implements json.Unmarshaler, allowing the genotype to be unmarshalled from JSON.
//
// TODO(Needs more validation)
func (g *NeatGenotype) UnmarshalJSON(bs []byte) error {
	mg := marshallableGenotype{}
	err := json.Unmarshal(bs, &mg)
	if err != nil {
		return err
	}
	g.neuronOrder = make([]NeatNeuronID, len(mg.Neurons))
	g.inverseNeuronOrder = make(map[NeatNeuronID]int)
	g.activations = make(map[NeatNeuronID]Activation)
	for ni, mn := range mg.Neurons {
		g.activations[mn.ID] = mn.Activation
		g.neuronOrder[ni] = mn.ID
		g.inverseNeuronOrder[mn.ID] = ni
	}
	g.weights = make(map[NeatSynapseID]float64)
	g.synapseEndpointLookup = make(map[NeatSynapseID]NeatSynapseEP)
	g.endpointSynapseLookup = make(map[NeatSynapseEP]NeatSynapseID)
	g.forwardSynapses = make([]NeatSynapseID, 0)
	g.backwardSynapses = make([]NeatSynapseID, 0)
	g.selfSynapses = make([]NeatSynapseID, 0)
	for _, ms := range mg.Synapses {
		ep := NeatSynapseEP{ms.From, ms.To}
		g.weights[ms.ID] = ms.Weight
		g.endpointSynapseLookup[ep] = ms.ID
		g.synapseEndpointLookup[ms.ID] = ep
		fromOrder := g.inverseNeuronOrder[ep.From]
		toOrder := g.inverseNeuronOrder[ep.To]
		if fromOrder < toOrder {
			g.forwardSynapses = append(g.forwardSynapses, ms.ID)
		} else if fromOrder > toOrder {
			g.backwardSynapses = append(g.backwardSynapses, ms.ID)
		} else {
			g.selfSynapses = append(g.selfSynapses, ms.ID)
		}
	}

	g.numInputs = mg.NumIn
	g.numOutputs = mg.NumOut
	g.maxSynapseValue = mg.MaxSynapseVal
	if err := g.Validate(); err != nil {
		return fmt.Errorf("genotype was invalid upon loading: %v", err)
	}
	return nil
}

// RenderDot returns a string in the DOT language that represents the genotype.
// This DOT code cannot be use to recreate the genotype, but can be used to visualise it using Graphviz.
func (g *NeatGenotype) RenderDot(width, height float64) string {
	graphDrawer := newSimpleGraphvizWriter()
	graphDrawer.writeGraphParam("rankdir", "LR")
	graphDrawer.writeGraphParam("ratio", "fill")
	graphDrawer.writeGraphParam("size", fmt.Sprintf("%v,%v", width, height))
	graphDrawer.writeGraphParam("layout", "dot")

	inputRanks := []string{}
	outputRanks := []string{}

	for no, nid := range g.neuronOrder {
		name := fmt.Sprintf("N%v", nid)
		label := fmt.Sprintf("N%v [%v]\n%v", nid, no, g.activations[nid])
		color := "black"
		if no < g.numInputs {
			color = "green"
			inputRanks = append(inputRanks, name)
		} else if no >= len(g.neuronOrder)-g.numOutputs {
			color = "red"
			outputRanks = append(outputRanks, name)
		}
		graphDrawer.writeNode(name, label, color)
	}

	graphDrawer.writeMinRank(inputRanks)
	graphDrawer.writeMaxRank(outputRanks)

	for wid, w := range g.weights {
		ep := g.synapseEndpointLookup[wid]
		of, ot := g.inverseNeuronOrder[ep.From], g.inverseNeuronOrder[ep.To]
		fromName := fmt.Sprintf("N%v", ep.From)
		toName := fmt.Sprintf("N%v", ep.To)
		label := fmt.Sprintf("C%v\n%.3f", wid, w)
		color := "black"
		if of >= ot {
			color = "red"
		}
		graphDrawer.writeEdge(fromName, toName, label, color)
	}
	return graphDrawer.dot()
}

// RenderImage returns an image of the genotype using graphviz.
func (g *NeatGenotype) RenderImage(width, height float64) image.Image {
	graph, err := graphviz.ParseBytes([]byte(g.RenderDot(width, height)))
	if err != nil {
		panic(fmt.Sprintf("error when creating a dot graph, this should not have happened (please report bug): %v", err))
	}
	gv := graphviz.New()
	img, err := gv.RenderImage(graph)
	if err != nil {
		panic(fmt.Sprintf("error when creating an image from dot, this should not have happened (please report bug): %v", err))
	}
	return img
}
