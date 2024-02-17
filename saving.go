package goevo

import (
	"encoding/json"
	"fmt"
)

// Make sure we implement json marshalling
var _ json.Marshaler = &Genotype{}
var _ json.Unmarshaler = &Genotype{}

type marshallableNeuron struct {
	ID         NeuronID   `json:"id"`
	Activation Activation `json:"activation"`
}

type marshallableSynapse struct {
	ID     SynapseID `json:"id"`
	From   NeuronID  `json:"from"`
	To     NeuronID  `json:"to"`
	Weight float64   `json:"weight"`
}

type marshallableGenotype struct {
	NumIn         int                   `json:"num_in"`
	NumOut        int                   `json:"num_out"`
	Neurons       []marshallableNeuron  `json:"neurons"`
	Synapses      []marshallableSynapse `json:"synapses"`
	MaxSynapseVal float64               `json:"max_synapse_val"`
}

// MarshalJSON implements json.Marshaler.
func (g *Genotype) MarshalJSON() ([]byte, error) {
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
func (g *Genotype) UnmarshalJSON(bs []byte) error {
	mg := marshallableGenotype{}
	err := json.Unmarshal(bs, &mg)
	if err != nil {
		return err
	}
	g.neuronOrder = make([]NeuronID, len(mg.Neurons))
	g.inverseNeuronOrder = make(map[NeuronID]int)
	g.activations = make(map[NeuronID]Activation)
	for ni, mn := range mg.Neurons {
		g.activations[mn.ID] = mn.Activation
		g.neuronOrder[ni] = mn.ID
		g.inverseNeuronOrder[mn.ID] = ni
	}
	g.weights = make(map[SynapseID]float64)
	g.synapseEndpointLookup = make(map[SynapseID]SynapseEP)
	g.endpointSynapseLookup = make(map[SynapseEP]SynapseID)
	g.forwardSynapses = make([]SynapseID, 0)
	g.backwardSynapses = make([]SynapseID, 0)
	g.selfSynapses = make([]SynapseID, 0)
	for _, ms := range mg.Synapses {
		ep := SynapseEP{ms.From, ms.To}
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
