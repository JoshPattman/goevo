package neat

import (
	"math"
	"math/rand"
	"slices"

	"github.com/JoshPattman/goevo"
)

// AddRandomNeuron adds a new neuron to the genotype on a random forward synapse.
// It will return false if there are no forward synapses to add to.
// The new neuron will have a random activation function from the given list of activations.
func AddRandomNeuron(g *Genotype, counter *goevo.Counter, activations ...goevo.Activation) bool {
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
func AddRandomSynapse(g *Genotype, counter *goevo.Counter, weightStd float64, recurrent bool) bool {
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
func MutateRandomSynapse(g *Genotype, std float64) bool {
	if len(g.weights) == 0 {
		return false
	}

	sid := randomMapKey(g.weights)
	g.weights[sid] = clamp(g.weights[sid]+rand.NormFloat64()*std, -g.maxSynapseValue, g.maxSynapseValue)

	return true
}

// RemoveRandomSynapse will remove a random synapse from the genotype.
// It will return false if there are no synapses to remove.
func RemoveRandomSynapse(g *Genotype) bool {
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
func ResetRandomSynapse(g *Genotype) bool {
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
func MutateRandomActivation(g *Genotype, activations ...goevo.Activation) bool {
	numHidden := len(g.neuronOrder) - g.numInputs - g.numOutputs
	if numHidden <= 0 {
		return false
	}
	i := g.numInputs + rand.Intn(numHidden)
	g.activations[g.neuronOrder[i]] = activations[rand.Intn(len(activations))]
	return true
}
