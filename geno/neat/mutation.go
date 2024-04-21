package neat

import "github.com/JoshPattman/goevo"

var _ goevo.MutationStrategy[*Genotype] = &StdMutation{}

// StdMutation is a reproduction strategy that uses a standard deviation for the number of mutations in each category.
// The standard deviation is not scaled by the size of the network, meaning that larger networks will tend to have more mutations than smaller networks.
type StdMutation struct {
	// The standard deviation for the number of new synapses
	StdNumNewSynapses float64
	// The standard deviation for the number of new recurrent synapses
	StdNumNewRecurrentSynapses float64
	// The standard deviation for the number of new neurons
	StdNumNewNeurons float64
	// The standard deviation for the number of synapses to mutate
	StdNumMutateSynapses float64
	// The standard deviation for the number of synapses to prune
	StdNumPruneSynapses float64
	// The standard deviation for the number of activations to mutate
	StdNumMutateActivations float64

	// The standard deviation for the weight of new synapses
	StdNewSynapseWeight float64
	// The standard deviation for the weight of mutated synapses
	StdMutateSynapseWeight float64

	// The maximum number of hidden neurons this mutation can add
	MaxHiddenNeurons int

	// The counter to use for new synapse IDs
	Counter *goevo.Counter
	// The possible activations to use for new neurons
	PossibleActivations []goevo.Activation
}

// Reproduce creates a new genotype by crossing over and mutating the given genotypes.
func (r *StdMutation) Mutate(g *Genotype) {
	for i := 0; i < stdN(r.StdNewSynapseWeight); i++ {
		AddRandomSynapse(g, r.Counter, r.StdNewSynapseWeight, false)
	}
	for i := 0; i < stdN(r.StdNumNewRecurrentSynapses); i++ {
		AddRandomSynapse(g, r.Counter, r.StdNewSynapseWeight, true)
	}
	for i := 0; i < stdN(r.StdNumNewNeurons); i++ {
		if r.MaxHiddenNeurons < 0 || g.NumHiddenNeurons() < r.MaxHiddenNeurons {
			AddRandomNeuron(g, r.Counter, r.PossibleActivations...)
		}
	}
	for i := 0; i < stdN(r.StdNumMutateSynapses); i++ {
		MutateRandomSynapse(g, r.StdMutateSynapseWeight)
	}
	for i := 0; i < stdN(r.StdNumPruneSynapses); i++ {
		RemoveRandomSynapse(g)
	}
	for i := 0; i < stdN(r.StdNumMutateActivations); i++ {
		MutateRandomActivation(g, r.PossibleActivations...)
	}
}
