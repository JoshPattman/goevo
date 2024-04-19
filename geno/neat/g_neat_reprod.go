package neat

import (
	"math/rand"

	"github.com/JoshPattman/goevo"
)

var _ goevo.Reproduction[*Genotype] = &StdReproduction{}
var _ goevo.Reproduction[*Genotype] = &ProbReproduction{}
var _ goevo.Reproduction[*Genotype] = &ScaledProbReproduction{}

// StdReproduction is a reproduction strategy that uses a standard deviation for the number of mutations in each category.
// The standard deviation is not scaled by the size of the network, meaning that larger networks will tend to have more mutations than smaller networks.
type StdReproduction struct {
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
func (r *StdReproduction) Reproduce(a, b *Genotype) *Genotype {
	g := a.CrossoverWith(b)

	for i := 0; i < stdN(r.StdNewSynapseWeight); i++ {
		g.AddRandomSynapse(r.Counter, r.StdNewSynapseWeight, false)
	}
	for i := 0; i < stdN(r.StdNumNewRecurrentSynapses); i++ {
		g.AddRandomSynapse(r.Counter, r.StdNewSynapseWeight, true)
	}
	for i := 0; i < stdN(r.StdNumNewNeurons); i++ {
		if r.MaxHiddenNeurons < 0 || g.NumHiddenNeurons() < r.MaxHiddenNeurons {
			g.AddRandomNeuron(r.Counter, r.PossibleActivations...)
		}
	}
	for i := 0; i < stdN(r.StdNumMutateSynapses); i++ {
		g.MutateRandomSynapse(r.StdMutateSynapseWeight)
	}
	for i := 0; i < stdN(r.StdNumPruneSynapses); i++ {
		g.RemoveRandomSynapse()
	}
	for i := 0; i < stdN(r.StdNumMutateActivations); i++ {
		g.MutateRandomActivation(r.PossibleActivations...)
	}

	return g
}

// ProbReproduction is a reproduction strategy that uses probabilities for the number of mutations in each category.
// The probabilities are not scaled by the size of the network, meaning that larger networks will tend to have the same number of mutations as smaller networks.
// ProbReproduction can also only perform one mutation per category.
type ProbReproduction struct {
	// The probability of adding a new neuron
	NewNeuronProbability float64
	// The probability of adding a new synapse
	NewSynapseProbability float64
	// The probability of adding a new recurrent synapse
	NewRecurrentSynapseProbability float64
	// The probability of removing a synapse
	RemoveSynapseProbability float64
	// The probability of mutating a synapse
	MutateSynapseProbability float64
	// The probability of mutating an activation
	MutateActivationProbability float64
	// The probability of setting a synapse to 0
	SetSynapseZeroProbability float64

	// The standard deviation for the weight of new synapses
	NewSynapseStd float64
	// The standard deviation for the weight of mutated synapses
	MutateSynapseStd float64
	// The possible activations to use for new neurons
	Activations []goevo.Activation

	// The probability of using the unfit parent as the base for the child
	UseUnfitParentProbability float64

	// The maximum number of hidden neurons this mutation can add
	MaxHiddenNeurons int

	// The counter to use for new synapse IDs
	Counter *goevo.Counter
}

// Reproduce creates a new genotype by crossing over and mutating the given genotypes.
func (r *ProbReproduction) Reproduce(a, b *Genotype) *Genotype {
	if rand.Float64() < r.UseUnfitParentProbability {
		a, b = b, a
	}
	g := a.CrossoverWith(b)

	if len(g.activations) < r.MaxHiddenNeurons && rand.Float64() < r.NewNeuronProbability {
		g.AddRandomNeuron(r.Counter, r.Activations...)
	}
	if rand.Float64() < r.NewSynapseProbability {
		g.AddRandomSynapse(r.Counter, r.NewSynapseStd, false)
	}
	if rand.Float64() < r.NewRecurrentSynapseProbability {
		g.AddRandomSynapse(r.Counter, r.NewSynapseStd, true)
	}
	if rand.Float64() < r.RemoveSynapseProbability {
		g.RemoveRandomSynapse()
	}
	if rand.Float64() < r.SetSynapseZeroProbability {
		g.ResetRandomSynapse()
	}
	if rand.Float64() < r.MutateSynapseProbability {
		g.MutateRandomSynapse(r.MutateSynapseStd)
	}
	if rand.Float64() < r.MutateActivationProbability {
		g.MutateRandomActivation(r.Activations...)
	}
	return g
}

// ScaledProbReproduction is a reproduction strategy that uses probabilities for the number of mutations in each category.
// The probabilities are scaled by the size of the network, meaning that larger networks will tend to have more mutations than smaller networks.
type ScaledProbReproduction struct {
	// Probability of creating a new neuron per synapse
	NewNeuronProbability float64
	// Probability of creating a synapse per (node^2)
	NewSynapseProbability float64
	// Probability of creating a recurrent synapse per (node^2)
	NewRecurrentSynapseProbability float64
	// Probability of removing each synapse
	RemoveSynapseProbability float64
	// Probability of mutating each synapse
	MutateSynapseProbability float64
	// Probability of mutating each activation
	MutateActivationProbability float64
	// Probability of setting each synapse to 0
	SetSynapseZeroProbability float64

	// Standard deviation for new synapse weights
	NewSynapseStd float64
	// Standard deviation for mutated synapse weights
	MutateSynapseStd float64
	// Possible activations to use for new neurons
	Activations []goevo.Activation

	// Probability of using the unfit parent as the base for the child
	UseUnfitParentProbability float64

	// Maximum number of hidden neurons this mutation can add
	MaxHiddenNeurons int

	// Counter to use for new synapse IDs
	Counter *goevo.Counter
}

// Reproduce creates a new genotype by crossing over and mutating the given genotypes.
func (r *ScaledProbReproduction) Reproduce(a, b *Genotype) *Genotype {
	if rand.Float64() < r.UseUnfitParentProbability {
		a, b = b, a
	}
	g := a.CrossoverWith(b)

	// Can only add a neuron on each forward synapse
	numNewNeuronPositions := len(g.forwardSynapses)
	for i := 0; i < numNewNeuronPositions; i++ {
		if len(g.activations) < r.MaxHiddenNeurons && rand.Float64() < r.NewNeuronProbability {
			g.AddRandomNeuron(r.Counter, r.Activations...)
		}
	}
	// Attempting to estimate number of new possible synapses
	// For now  we will approximate this to be numnodes^2
	numNewSynapsePositions := len(g.activations) * len(g.activations)
	for i := 0; i < numNewSynapsePositions; i++ {
		if rand.Float64() < r.NewSynapseProbability {
			g.AddRandomSynapse(r.Counter, r.NewSynapseStd, false)
		}
		if rand.Float64() < r.NewRecurrentSynapseProbability {
			g.AddRandomSynapse(r.Counter, r.NewSynapseStd, true)
		}
	}
	// Can remove any synapse
	numRemovableSynapses := len(g.weights)
	for i := 0; i < numRemovableSynapses; i++ {
		if rand.Float64() < r.RemoveSynapseProbability {
			g.RemoveRandomSynapse()
		}
	}
	// Can zero any synapse
	numZeroableSynapses := len(g.weights)
	for i := 0; i < numZeroableSynapses; i++ {
		if rand.Float64() < r.SetSynapseZeroProbability {
			g.ResetRandomSynapse()
		}
	}
	// Can mutate any synapse
	numMutateableSynapses := len(g.weights)
	for i := 0; i < numMutateableSynapses; i++ {
		if rand.Float64() < r.MutateSynapseProbability {
			g.MutateRandomSynapse(r.MutateSynapseStd)
		}
	}
	// Can mutate any activation
	numMutateableActivations := len(g.activations)
	for i := 0; i < numMutateableActivations; i++ {
		if rand.Float64() < r.MutateActivationProbability {
			g.MutateRandomActivation(r.Activations...)
		}
	}
	return g
}
