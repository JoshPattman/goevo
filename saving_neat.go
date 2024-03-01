package goevo

import (
	"encoding/json"
	"fmt"
)

// Make sure we implement json marshalling
var _ json.Marshaler = &NEATGenotype{}
var _ json.Unmarshaler = &NEATGenotype{}

type marshallableNeuron struct {
	ID         NEATNeuronID `json:"id"`
	Activation Activation   `json:"activation"`
}

type marshallableSynapse struct {
	ID     NEATSynapseID `json:"id"`
	From   NEATNeuronID  `json:"from"`
	To     NEATNeuronID  `json:"to"`
	Weight float64       `json:"weight"`
}

type marshallableGenotype struct {
	NumIn         int                   `json:"num_in"`
	NumOut        int                   `json:"num_out"`
	Neurons       []marshallableNeuron  `json:"neurons"`
	Synapses      []marshallableSynapse `json:"synapses"`
	MaxSynapseVal float64               `json:"max_synapse_val"`
}

// MarshalJSON implements json.Marshaler.
func (g *NEATGenotype) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON implements json.Unmarshaler.
// TODO: needs more validation
func (g *NEATGenotype) UnmarshalJSON(bs []byte) error {
	mg := marshallableGenotype{}
	err := json.Unmarshal(bs, &mg)
	if err != nil {
		return err
	}
	g.neuronOrder = make([]NEATNeuronID, len(mg.Neurons))
	g.inverseNeuronOrder = make(map[NEATNeuronID]int)
	g.activations = make(map[NEATNeuronID]Activation)
	for ni, mn := range mg.Neurons {
		g.activations[mn.ID] = mn.Activation
		g.neuronOrder[ni] = mn.ID
		g.inverseNeuronOrder[mn.ID] = ni
	}
	g.weights = make(map[NEATSynapseID]float64)
	g.synapseEndpointLookup = make(map[NEATSynapseID]NEATSynapseEP)
	g.endpointSynapseLookup = make(map[NEATSynapseEP]NEATSynapseID)
	g.forwardSynapses = make([]NEATSynapseID, 0)
	g.backwardSynapses = make([]NEATSynapseID, 0)
	g.selfSynapses = make([]NEATSynapseID, 0)
	for _, ms := range mg.Synapses {
		ep := NEATSynapseEP{ms.From, ms.To}
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
