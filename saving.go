package goevo

import "encoding/json"

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
	NumIn    int                   `json:"num_in"`
	NumOut   int                   `json:"num_out"`
	Neurons  []marshallableNeuron  `json:"neurons"`
	Synapses []marshallableSynapse `json:"synapses"`
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
	mg := marshallableGenotype{g.numInputs, g.numOutputs, mns, mss}
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
	for _, ms := range mg.Synapses {
		ep := SynapseEP{ms.From, ms.To}
		g.weights[ms.ID] = ms.Weight
		g.endpointSynapseLookup[ep] = ms.ID
		g.synapseEndpointLookup[ms.ID] = ep
	}
	return nil
}
